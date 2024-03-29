 # SLIP-Flow Notes

- next

 - actors
  + flow-read-file-actor
  + flow-read-json-actor
  + flow-read-csv-actor
  - flow-read-directory-actor
  - flow-write-file-actor
  - flow-delete-file-actor
  - flow-read-xml-actor
   - follow slip/pkg/xml format but as json

  - flow-watch-directory-actor
  - flow-tail-file-actor

 - add http server actor
  - http server could keep a reply channel
   - limitation is no external (other process) terminations (not really a restriction)

----------------
- flow :svg

------------------
  - trigger tasks/actors
   - are they needed or better to call from outside flow?
   - http-server-actor
   - queue-actor input actor for trigger tasks
   - schedule-actor
   - for
    - keeps servers in the flow (not sure if this is a good thing though)
   - against
    - requires loop back to original to respond for an http server
     - multiple flow entry points since nothing in box identifies as a second submission
     - no longer just flow forward so sub-flows need to know about tasks in other flows
    - allows server creation and configuration outside the flow so no need to deal with task/actor configuration restrictions


---------------
- text/x-common-lisp
- text/x-emacs-lisp
- application/lisp (not a recognized content-type)
- text/lisp (not a recognized content-type)
