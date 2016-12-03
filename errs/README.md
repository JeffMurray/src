package errs formalizes the seperation of system information and user information when handeling errors durng code execution.

The default Error() funcs do not give away state informatiion, but do give token trails to locate the errror in code.

Use go to compile the newid project to get a unique token to hardcode for each new error.

Source in the ezjson package gives a good example of using the errs package.

Assigning a unique token to every new error, whether system or user, makes locating where the problem happened in code a grep away.
