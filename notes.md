# SLIP-Flow Notes

- next

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


- dev env (in slip) or maybe emacs integration is enough
 - repl
 - ability to save functions to a file once they are good
  - maybe save world but just for changes
   - world.lisp
  - would need to remember function code
  - maybe list all defs (defvar, defun, defflavor, defmethod, defconstant)
   - pick which ones or all to write to a file
  - need an edit-function with option for external editor or in repl

- text/x-common-lisp
- text/x-emacs-lisp
- application/lisp (not a recognized content-type)
- text/lisp (not a recognized content-type)
