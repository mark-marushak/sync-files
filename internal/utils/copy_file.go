package utils

import cp "github.com/otiai10/copy"

func CopyFile(src, dest string) error {
	//sourceFile, err := os.Open(src)
	//if err != nil {
	//	return err
	//}
	//defer sourceFile.Close()
	//
	//destFile, err := os.Create(dst)
	//if err != nil {
	//	return err
	//}
	//defer destFile.Close()
	//
	//_, err = io.Copy(destFile, sourceFile)
	//if err != nil {
	//	return err
	//}
	//
	//err = destFile.Sync()
	//return err

	opt := cp.Options{
		Sync: true,
	}

	err := cp.Copy(src, dest, opt)
	if err != nil {
		return err
	}

	return nil
}
