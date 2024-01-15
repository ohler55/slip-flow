# SLIP-Flow Notes

- next
 - add options for write in log-error-actor

 - actors
  + flow-exit-actor (places box on exit-channel)
  + log-error-actor (logs error or what ever is in box if not an error)
  - inspect-actor [prints out box, for debugging]
  - queue-actor input actor for trigger tasks
  - split-actor [multiple transitions in parallel]
  - merge-actor [merge back multiple branch after a split]
  - sub-flow-actor
  - http-server-actor
  - http-get-actor
  - http-post-actor
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
