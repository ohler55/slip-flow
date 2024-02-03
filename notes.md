# SLIP-Flow Notes

- next

  - trigger tasks/actors
   - are they needed or better to call from outside flow?
   - http-server-actor
   - queue-actor input actor for trigger tasks
   - schedule-actor

- flow-editor
 - should be part of task and flow
  - keeping separate will be hard to keep in sync
 - allow tasks and flow to be updated after creation
 - provide write for a flow to write and file that can be used to recreate
 - task changes or updates (attributes and get and set)
  - put each in task-x.go, task-y.go, task-svg.go
  - :svg [string]
  - :x [fixnum]
  - :y
 - flow
  - :width
  - :height
  - :icon-width
  - :icon-height
  - :background [string]
 - flow :link
  - :points
 - optional method for actors to allow creation
  - init-key-values => (:foo 2 :bar "xyz")
   - what keywords and values are needed to create the same actor instance

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
