package babe

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/mrnavastar/assist/bytes"
	fss "github.com/mrnavastar/assist/fs"
	"golang.org/x/sync/errgroup"
)

type JarMember struct {
	Name   string
	Buffer *bytes.Buffer
	delete bool
}

func JarMemberFromFile(filename string) (member JarMember, err error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return member, err
	}
	member.Name = filename
	member.Buffer = &bytes.Buffer{Data: &data, Index: 0}
	return member, nil
}

func JarMemberFromString(name string, str string) (member JarMember) {
	data := []byte(str)
	member.Name = name
	member.Buffer = &bytes.Buffer{Data: &data, Index: 0}
	return member
}

func (member *JarMember) Delete() {
	member.delete = true
}

func (member *JarMember) GetAsClass() (Class, error) {
	if !strings.HasSuffix(member.Name, ".class") {
		return Class{}, ErrNotClass
	}
	var class Class
	if err := class.Read(*member.Buffer.Data); err != nil {
		return class, err
	}
	return class, nil
}

type Jar struct {
	Name  string
	c     chan *JarMember
	tasks *errgroup.Group
	group *errgroup.Group
}

func (jar *Jar) Task(task func(jar *Jar) error) {
	jar.tasks.Go(func() error {
		return task(jar)
	})
}

func (jar *Jar) Add(member JarMember) {
	jar.c <- &member
}

func (jar *Jar) Wait() error {
	if err := jar.tasks.Wait(); err != nil {
		return err
	}
	close(jar.c)
	return jar.group.Wait()
}

func ForJarMember(filename string, iter func(*JarMember) error) error {
	reader, err := zip.OpenReader(filename)
	if err != nil {
		return err
	}
	defer reader.Close()

	errs, _ := errgroup.WithContext(context.Background())
	for _, file := range reader.File {
		errs.Go(func() error {
			if file.FileInfo().IsDir() {
				return nil
			}

			f, err := file.Open()
			if err != nil {
				return err
			}

			member := JarMember{file.Name, bytes.NewBuffer(), false}
			if _, err = io.Copy(member.Buffer, f); err != nil {
				return err
			}
			f.Close()

			if err = iter(&member); err != nil {
				return err
			}
			return nil
		})
	}
	return errs.Wait()
}

func CreateJar(filename string) (jar Jar) {
	jar.Name = path.Base(filename)
	jar.c = make(chan *JarMember)
	jar.tasks, _ = errgroup.WithContext(context.Background())
	jar.group, _ = errgroup.WithContext(context.Background())

	jar.group.Go(func() error {
		file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return err
		}
		writer := zip.NewWriter(file)

		for {
			member, ok := <-jar.c
			if !ok {
				break
			}

			w, err := writer.CreateHeader(&zip.FileHeader{Name: member.Name, Method: zip.Deflate})
			if err != nil {
				return err
			}
			_, err = w.Write(*member.Buffer.Data)
			if err != nil {
				return err
			}
		}
		return writer.Close()
	})
	return jar
}

func ModifyJar(filename string, modifier func(*JarMember) error) error {
	if !fss.Exists(filename) {
		return fmt.Errorf("%s does not exist", filename)
	}

	jar := CreateJar(filename + "-modified.zip")
	jar.Task(func(jar *Jar) error {
		return ForJarMember(filename, func(member *JarMember) error {
			if err := modifier(member); err != nil {
				return err
			}

			if !member.delete {
				jar.Add(*member)
			}
			return nil
		})
	})

	if err := os.Rename(filename+"-modified.zip", filename); err != nil {
		return err
	}
	return jar.Wait()
}
