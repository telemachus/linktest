# linktest: Test the links in one or more files

`linktest` searches for link rot in HTML files.  By default, it only reports
the status of links that don't return 200.  If you want to see the status of
all links, add `-verbose`.

`linktest` has two important limitations.  First, it does not handle `stdin`
only files.  Second, expects the files to be HTML.

Finally, `linktest` is very new and probably closer to alpha than beta.
