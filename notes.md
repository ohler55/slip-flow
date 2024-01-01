# SLIP-Flow Notes

- next

 - new flow-flavor branch
  - add flavor to flow.go
   - docs
    - describe as channel where submit or trigger starts and exit results on channel
    - describe error task name is special
    - will need an exit actor that puts box on exit-channel

  - methods
   - :init
    - option to set manager which pulls in logger
    - option to set logger directly
    - option for a completed-channel or exit-channel or finished-channel or result-channel
   - :add-task (task keywords and values or just task?)
    - if args then no worries about adding the same task twice
    - use taskInitCaller?
   - :link (name from to) [all args are strings)
   - :unlink (task-name link-name)
   - :remove-task (task-name)
   - :start
   - :shutdown
   - :tasks
   - :find-task (task-name)
   - :entry
   - :set-entry (task-name)
   - :submit (box)

  - test task :receive panics in actor
   - make sure string and Stringer both work as well as error
    - for string and error need go actor


   - task struct
    - name
    - self points back to task instance
    - links map[string]*Link
    - actors []*flavors.Instance
    - function (if using a function or lambda)
    - queue chan box instance
    - workers (used to put nil on chan to stop)



- actors
 - queue input actor for trigger tasks
 - splitter and merger
 - sub-flow
 - ...

 - classes/flavors
  - flow-manager-flavor
   - :init [directory of flows or config file or config args]
   - :load [read/load a config file then add]
   - :add [from a bag or lisp config]
   - :find [get a flow by name]
   - :flows [all flows, maybe with pattern to match]
   - :remove
   - :logger [return gi:logger of the manager]
