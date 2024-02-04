# SLIP-Flow Notes

- next

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

- writeable
 - flow
  - :width
  - :height
  - :icon-width
  - :icon-height
  - :background [string]
  - :link
   - :points keyword as list of x y pairs ((1 1)(10 20))
  - :write
   - writes code to create in a let or let*
   - use (read x) to load or add to a group
    - add read to slip
  - :draw
   - returns a string that is an SVG
 - task
  - :svg [string] [task-svg.go]
  - :x [fixnum] [task-x.go]
  - :y [fixnum] [task-y.go]
 - actors
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
