# SLIP-Flow Notes

- next
 - change examples to use flow to create flow, add task
 - change tests to use MakeTask() or flow then add-task

  - methods
   + :init
   + :name
   + :exit-channel
   + :start
   + :shutdown
   + :running
   + :add-task
   + :tasks
   + :remove-task (task-name)
   + :find-task (task-name)
   + :entry
   + :set-entry (task-name)
   - :link (name from to) [all args are strings)
   - :unlink (task-name link-name)
   - :submit (box)
   - :metrics
    - task entry with empty track increments flow received
    - any termination task (ExitTask or flag on task to indicate it is last?) marks as processed and updates duration
     - maybe just any task with no outgoing links
      - actor returns box and link so box and nil link-name with no links

  - task :links
   - return assoc with link name and destination task

  - test task :receive panics in actor
   - make sure string and Stringer both work as well as error
    - for string and error need go actor

 - flow-manager

 - examples
  - multiple task flows

 - actors
  - exit actor (places box on exit-channel)
  - error-logger (logs error or what ever is in box if not an error)
  - queue input actor for trigger tasks
  - splitter and merger
  - sub-flow
  - ...

- UI
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
 - update design.md
 - optional method for actors to allow creation
  - init-key-values => (:foo 2 :bar "xyz")
   - what keywords and values are needed to create the same actor instance


 - classes/flavors
  - flow-manager-flavor
   - :init [directory of flows or config file or config args]
   - :load [read/load a config file then add]
   - :add [from a bag or lisp config]
   - :find [get a flow by name]
   - :flows [all flows, maybe with pattern to match]
   - :remove
   - :logger [return gi:logger of the manager]
