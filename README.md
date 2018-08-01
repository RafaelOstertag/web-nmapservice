Nmap web service.

Configuration
---

The service is configured via environment variables.

* `NMAP_CMD`: file path of the nmap executable. Defaults to `nmap`.
* `LISTEN`: address and port to listen on. This value is passed to `http.ListenAndServe()`. Defaults to `:8081`.
* `LOGPATH`: path to log directory. When set, logs are written to this directory.
* `CONSUL`: host and port to consul. Defaults to `gizmo.kruemel.home:8500`