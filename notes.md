# SLIP-Flow Notes

- next
 - flow-box-flavor
 - has-logger-flavor (abstract flavor)

- design
 - flow package in separate repo
  - design/model
   - immutable bag or similar
   - separate routine for each node/task
    - or option for separate thread
     - **maybe option for number of workers**
      - 0 means inline, > 0 means work-queue (channel) and workers
  - flow-manager
   - logger
   - map of flows by name
  - flow
   - name
   - log level
   - entry
   - tasks
  - task or node
   - name
   - links map[string]*Link
   - actors []*flavors.Instance
    - :receive
     - calls :transition on parent task to move to linked task
   - function (if using a function or lambda)
   - queue chan *Box
   - worker-count
   - workers (if actor is an instance)
    - instances in the worker loops if instances
    - else just use the function
  - link
   - name
   - target
  - box
   - error [var]
	flavor.DefMethod(":init", "", initCaller{}) [provide tracking id or track from other box]
	flavor.DefMethod(":set", "", setCaller{})
	flavor.DefMethod(":parse", "", parseCaller{})
	flavor.DefMethod(":read", "", readCaller{})
	flavor.DefMethod(":get", "", getCaller{})
	flavor.DefMethod(":has", "", hasCaller{})
	flavor.DefMethod(":remove", "", removeCaller{})
	flavor.DefMethod(":modify", "", modifyCaller{})
	flavor.DefMethod(":native", "", nativeCaller{})
	flavor.DefMethod(":write", "", writeCaller{})
	flavor.DefMethod(":walk", "", walkCaller{})
	flavor.DefMethod(":bag", "", bagCaller{})
	flavor.DefMethod(":native", "", nativeCaller{})
	flavor.DefMethod(":tracking-id", "", trackingIDCaller{})
	flavor.DefMethod(":freeze", "", freezeCaller{})
	flavor.DefMethod(":track", "", trackCaller{}) [instance]
	flavor.DefMethod(":events", "", eventsCaller{})
	flavor.DefMethod(":scan", "", scanCaller{}) [adds and entry to the track]

  - config format (lisp or json)
   - flow
    - name
    - log-level
    - entry (string)
    - tasks
     - name
     - log-level
     - links (names and targets)
     - function (optional)
     - actor
      - flavor
      - init key/values
     - worker-count

   - json format in a bag with option for lisp

- actors
 - queue input actor for trigger tasks
 - splitter and merger
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
  - has-logger-flavor
   - logger [read only and set by container or on create]
   - :log-level
   - :set-log-level
   - :error
   - :warn
   - :info
   - :debug
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

  - flow-link-flavor or flow-transition-flavor

  - flow-task-flavor (has-logger)
   - :transition (box &optional wait)
   - :flow [containing flow]
   - actor [var]
   - function [var]
   -
  - flow-box-flavor
   - model after bag but add tracking-id
   - keep flag to indicate if it is immutable
    - set flag on call to transition or receive
     - only need to set the flag if async (0 < workers or has a queue)

   - task-flavor
    - methods
     - start
      - starts processing loop
     - stop
     - submit box/data/bag
      - drops data on to processing channel
      - initially copy but later wrap with box that dups on set
       - or maybe enhance bag to have option for copy on set (immuttable flag)
     - handle-error
     - flow return flow task is in
     - links - returns link names with task as assoc list
    - subclass for specific behavior
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
