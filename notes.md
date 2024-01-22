# SLIP-Flow Notes

- next
 - actors
  + flow-exit-actor (places box on exit-channel)
  + log-error-actor (logs error or what ever is in box if not an error)
  + inspect-actor [prints out box, for debugging]
  + jump-actor
  + split-actor
  - merge-actor [merge back multiple branch after a split]
   - number of boxes to receive for an id
   - timeout if not all received in alloted time
   - keep a mutex protected map
   - need a timeout check loop (check every timeout/2 ?)
   - merge box track and sort by time
   - box content merged
    - latest overrides if necessary
    - jp.Walk non-primary, check values and add if missing, ignore if different
  - http-client-actor (specify get, post, etc)
  - http-get-actor - is this needed?
  - http-post-actor - is this needed?
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
