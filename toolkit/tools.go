package toolkit

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const randStringSource = "abcdefghijklmnopqrstuvwxyzACDEFGHIJKLMNOPQRSTUVWXYZ0123456789_+"

type Tools struct {
	MaxFileSize      int
	AllowedFileTypes []string
}

func (t *Tools) RandomString(n int) string {
	s, r := make([]rune, n), []rune(randStringSource)
	for i := range s {
		p, _ := rand.Prime(rand.Reader, len(r))
		x, y := p.Uint64(), uint64(len(r))
		s[i] = r[x%y]
	}
	return string(s)
}

// For infoi abt uploaded file
type UploadedFile struct {
	NewFileName      string
	OriginalFileName string
	FileSize         int64
}

func (t *Tools) UploadFiles(r *http.Request, uploadDir string, rename ...bool) ([]*UploadedFile, error) {
	renameFile := true
	if len(rename) > 0 {
		renameFile = rename[0]
	}

	var uploadedFiles []*UploadedFile
	if t.MaxFileSize == 0 {
		t.MaxFileSize = 1024 * 1024 * 1024
	}
	err := r.ParseMultipartForm(int64(t.MaxFileSize))
	if err != nil {
		return nil, errors.New("the uploaded file is too big")
	}
	for _, fHeaders := range r.MultipartForm.File {
		for _, hdr := range fHeaders {
			uploadedFiles, err = func(uploadedFiles []*UploadedFile) ([]*UploadedFile, error) {
				var uploadedFile UploadedFile
				infile, err := hdr.Open()
				if err != nil {
					return nil, err
				}
				defer infile.Close()

				buff := make([]byte, 512)
				_, err = infile.Read(buff)
				if err != nil {
					return nil, err
				}

				allowed := false
				fileType := http.DetectContentType(buff)

				if len(t.AllowedFileTypes) > 0 {
					for _, x := range t.AllowedFileTypes {
						if strings.EqualFold(fileType, x) {
							allowed = true
						}
					}
				} else {
					allowed = true
				}

				if !allowed {
					return nil, errors.New("the uploaded file type is not permitted")
				}

				_, err = infile.Seek(0, 0)
				if err != nil {
					return nil, err
				}
				if renameFile {
					uploadedFile.NewFileName = fmt.Sprintf("%s%s", t.RandomString(25), filepath.Ext(hdr.Filename))
				} else {
					uploadedFile.NewFileName = hdr.Filename
				}
				var outfile *os.File
				defer outfile.Close()

				if outfile, err = os.Create(filepath.Join(uploadDir, uploadedFile.NewFileName)); err != nil {
					return nil, err
				} else {
					fileSize, err := io.Copy(outfile, infile)
					if err != nil {
						return nil, err
					}
					uploadedFile.FileSize = fileSize
				}

				uploadedFiles = append(uploadedFiles, &uploadedFile)

				return uploadedFiles, nil
			}(uploadedFiles)
			if err != nil {
				return uploadedFiles, err
			}
		}
	}
	return uploadedFiles, nil
}

// func (t *Tools) UploadFiles(r *http.Request, uploadDir string, rename ...bool) ([]*UploadedFile, error) {
// 	renameFile := true
// 	if len(rename) > 0 {
// 		renameFile = rename[0]
// 	}

// 	var uploadedFiles []*UploadedFile
// 	if t.MaxFileSize == 0 {
// 		t.MaxFileSize = 1024 * 1024 * 1024
// 	}
// 	err := r.ParseMultipartForm(int64(t.MaxFileSize))
// 	if err != nil {
// 		return nil, errors.New("The uploaded file is too big!")
// 	}
// 	for _, fHeaders := range r.MultipartForm.File {
// 		for _, hdr := range fHeaders {
// 			uploadedFiles, err = func(uploadedFiles []*UploadedFile) ([]*UploadedFile, error) {
// 				var uploadedFile UploadedFile
// 				infile, err := hdr.Open()
// 				if err != nil {
// 					return nil, err
// 				}
// 				defer infile.Close()
// 				buff := make([]byte, 512)
// 				_, err = infile.Read(buff)
// 				if err != nil {
// 					return nil, err
// 				}
// 				// TODO: check to see if the filetype is permitted
// 				allowed := false
// 				fileType := http.DetectContentType(buff)
// 				// allowedTypes := []string{
// 				// 	"image/jpeg", "image/png", "image/gif",
// 				// }
// 				if len(t.AllowedFileTypes) > 0 {
// 					for _, x := range t.AllowedFileTypes {
// 						if strings.EqualFold(fileType, x) {
// 							allowed = true
// 						}
// 					}
// 				} else {
// 					allowed = true
// 				}
// 				if !allowed {
// 					return nil, errors.New("the uploaded file type is not permitted")
// 				}
// 				// Go back to the beginning of file
// 				_, err = infile.Seek(0, 0)
// 				if err != nil {
// 					return nil, err
// 				}
// 				if renameFile {
// 					uploadedFile.NewFileName = fmt.Sprintf("%s%s", t.RandomString(25), filepath.Ext(hdr.Filename))
// 				} else {
// 					uploadedFile.NewFileName = hdr.Filename

// 				}
// 				var outfile *os.File
// 				defer outfile.Close()

// 				outfile, err = os.Create(filepath.Join(uploadDir, uploadedFile.NewFileName))
// 				if err != nil {
// 					return nil, err
// 				}
// 				fileSize, err := io.Copy(outfile, infile)
// 				if err != nil {
// 					return nil, err
// 				}
// 				uploadedFile.FileSize = fileSize
// 				uploadedFiles = append(uploadedFiles, &uploadedFile)
// 				return uploadedFiles, nil

// 			}(uploadedFiles)
// 			if err != nil {
// 				return uploadedFiles, err
// 			}
// 		}
// 	}
// 	return uploadedFiles, nil
// }
