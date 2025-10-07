package pipeline

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"

	"github.com/KhoalaS/godel/pkg/utils"
	"github.com/mholt/archives"
)

func CreateUnrarNode() Node {
	return Node{
		Type:     "unrar",
		Run:      UnrarNodeFunc,
		Name:     "Unrar",
		Status:   StatusPending,
		Category: NodeCategoryUtility,
		Io: map[string]*NodeIO{
			"file": {
				Id:        "file",
				ValueType: ValueTypeFile,
				Label:     "File",
				Required:  true,
				Type:      IOTypeConnectedOnly,
			},
			"password": {
				Id:        "password",
				ValueType: ValueTypeString,
				Label:     "Password",
				Required:  false,
				Value:     "",
				Type:      IOTypeInput,
			},
		},
	}
}

func UnrarNodeFunc(ctx context.Context, node Node, pipeline IPipeline) error {
	unrarExists, _ := utils.ExecutableExists("unrar")
	file, ok := node.Io["file"].Value.(IFile)
	if !ok {
		return errors.New("missing file input")
	}

	absolutePath, err := file.GetAbsolutePath()
	if err != nil {
		return err
	}

	destFolder := file.GetDestinationFolder()

	if !unrarExists {
		fileHandle, err := file.GetFileHandle()
		if err != nil {
			return err
		}

		format, stream, err := archives.Identify(ctx, absolutePath, fileHandle)
		if err != nil {
			return err
		}

		if ex, ok := format.(archives.Extractor); ok {
			err = ex.Extract(ctx, stream, func(ctx context.Context, info archives.FileInfo) error {
				if info.FileInfo.IsDir() {
					os.Mkdir(filepath.Join(destFolder, info.NameInArchive), 0755)
					return nil
				} else {
					outdir := filepath.Dir(info.NameInArchive)
					err := os.MkdirAll(filepath.Join(destFolder, outdir), 0755)
					if err != nil {
						return err
					}

					infile, err := info.Open()
					if err != nil {
						fmt.Println("1")
						return err
					}

					outfile, err := os.Create(filepath.Join(destFolder, info.NameInArchive))
					if err != nil {
						fmt.Println("2")
						return err
					}

					io.Copy(outfile, infile)
					outfile.Close()
					infile.Close()
				}

				return nil
			})

			if err != nil {
				debug.PrintStack()
				fmt.Println("3")
				return err
			}
		} else {
			return errors.New("file can not be extracted")
		}

		fileHandle.Close()
		return nil
	}

	unrarCommand := exec.Command("unrar", "x")
	if password, ok := node.Io["password"].Value.(string); ok && password != "" {
		unrarCommand.Args = append(unrarCommand.Args, fmt.Sprintf("-p%s", password))
	}

	unrarCommand.Args = append(unrarCommand.Args, absolutePath, file.GetDestinationFolder())
	err = unrarCommand.Start()
	if err != nil {
		return err
	}

	err = unrarCommand.Wait()
	if err != nil {
		return err
	}

	return nil
}
