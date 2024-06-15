package file

import "os"

func WriteFile(f *File) error {
	// assume that file name is the absolute path
	outFile, err := os.Create(f.Name)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = outFile.Write(f.Data)
	if err != nil {
		return err
	}

	return nil
}
