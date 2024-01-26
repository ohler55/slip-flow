# SLIP-Flow Notes

- next
 - actors
  + flow-exit-actor (places box on exit-channel)
  + log-error-actor (logs error or what ever is in box if not an error)
  + inspect-actor [prints out box, for debugging]
  + jump-actor
  + split-actor
  + merge-actor

  - http-client-actor (specify get, post, etc)
   - init
    - method
    - host
    - port
    - url
    - proto
    - header
    - trailer
    - body (string or stream)
    - timeout
    - response-path
     - place in box at path
     - encode header and other response data
      - if content type is json then parse and store
   - :perform
    - pull params from box to override init params
     - use top path and fixed sub paths?
     - options for each init param can be a path instead
      - how to know if it is a path
       - since actors are created with code paths can be bag-paths
      - or maybe a function/lambda to generate the value
    - log info for request and response

  - trigger tasks
   - are they needed or better to call from outside flow?
   - http-server-actor
   - queue-actor input actor for trigger tasks
   - schedule-actor

- flow-editor
 - should be part of task and flow
  - keeping separate will be hard to keep in sync
  - task
   - svg [string]
   - x [fixnum]
   - y
  - flow
   - width
   - height
   - icon-width
   - icon-height
   - background [string]
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
