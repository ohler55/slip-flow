# SLIP-Flow Notes

- next
 - create flow-task-flavor and :init
  - test with other methods

  - flow-task-flavor
   - initializers
    - :name
    - :actor - one or a list
     - function, instance, or list of instances
    - :worker (count)
     - if 0 then sync
     - should match actors but round-robin for assignment to loops if not
    - :logger (from can-log-flavor)
    - :log-level (from can-log-flavor)
   - methods
    - :receive
    - :transition
    - :start
    - :shutdown (&optional wait)
    - all from can-log-flavor

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
  - flow-flavor (has-logger-flavor)
   - :init [should take a config but allow for changes later]
   - :start [starts all tasks]
   - :stop &optional wait [all tasks]
   - :submit (box &optional wait) [or call it receive to match tasks]
   - :handle-error [called by tasks]
   - :manager
   - :tasks
   - :add-task
   - :remove-task
   - :find-task
   - :entry
   - :link (source link-name target) [source and target can be name or task itself or another flow]
    - maybe both target and source must be names to avoid cross linking
    - still allow target to be a flow though
   - :unlink (source link-name)
   - :set-entry (task-name)

  - flow-flavor
   - init should take a config but allow for changes later
   - can subclass for specific flows
   - methods
     - start (starts all tasks)
     - stop &optional wait (all tasks)
    - submit data &optional wait
    - handle-error
    - logger field points to gi/logger
    - tasks
    - add-task
    - set-entry
    - link name task (or task-name)
     - optional flow for external links
   - vars
    - entry
    - tasks
  - syntax for describing, json or lisp
