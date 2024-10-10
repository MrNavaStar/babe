package babe

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
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

func CreateJar(filename string) (chan *JarMember, *errgroup.Group) {
	c := make(chan *JarMember)

	errs, _ := errgroup.WithContext(context.Background())
	errs.Go(func() error {
		file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return err
		}
		writer := zip.NewWriter(file)

		for {
			member, ok := <-c
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
	return c, errs
} 

func ModifyJar(filename string, modifier func(*JarMember) error) error {
	if !fss.Exists(filename) {
		return fmt.Errorf("%s does not exist", filename)
	}

	c, errs := CreateJar(filename+"-modified.zip")

	err := ForJarMember(filename, func(member *JarMember) error {
		if err := modifier(member); err != nil {
			return err
		}

		if !member.delete {
			c <- member
		}
		return nil
	})
	close(c)

	if err != nil {
		return err
	}
	
	if err = os.Remove(filename); err != nil {
		return err
	}
	if err = os.Rename(filename+"-modified.zip", filename); err != nil {
		return err
	}

	return errs.Wait()
}
