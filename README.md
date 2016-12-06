package at_rest has some cross platform, non-expert friendly functions to encrypt and decrypt byte arrays, and small configuration files.

package ezjson makes working with newly decoded map[string]interface{} types almost as easy as javascript.

package errs a powerful way seperate sysetm from client/user errors, and tracing the path they took in code.

newid creates newid.exe that spits out a 6 char unique string chars a-z amd 0-9, handy for the error handeling scheme in package errs

package safe_file_relay implements the WriteAt and Read interfaces.  It is designed to concurently connect non-sequential WrtieAts from a sftp server to sequential Reads from an s3 bucket uploader.