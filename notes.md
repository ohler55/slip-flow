# SLIP-Flow Notes

- next

 - flow-group-flavor inherits can-log-flavor
  - :init
  - :add [flow-group-add]
  - :find [get a flow by name]
  - :flows [all flows, maybe with pattern to match]
  - :remove
  - :set-level :after
  - :start
  - :shutdown
  - add group to flow so sub-flow-actor can find sub-flow

 - actors
  + flow-exit-actor (places box on exit-channel)
  - log-error-actor (logs error or what ever is in box if not an error)
   - option to exit or continue flow
    - maybe just look at links. if no links then exit else follow any link
  - queue-actor input actor for trigger tasks
  - split-actor [multiple transitions in parallel]
  - merge-actor [merge back multiple branch after a split]
  - sub-flow-actor
  - http-server-actor
  - http-get-actor
  - http-post-actor
  - schedule-actor
  - inspect-actor [prints out box, for debugging]
  - file-monitor-actor [load and process file when it changes, or if dir when new file added]
  - file-saver-actor [write or append to file, name generated from content]

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
