# linktest: Test external links in HTML files

`linktest` searches for link rot in HTML files.  By default, it only reports
the status of links that don't return 200.  If you want to see the status of
all links, add `-verbose`.

`linktest` has two deliberate limitations.  First, it only handles files—not
`stdin`.  Second, it expects the files to be HTML.

Finally, `linktest` is very new and closer to alpha than beta.
