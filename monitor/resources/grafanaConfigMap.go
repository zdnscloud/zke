package resources

const GrafanaConfigMapYaml = `
apiVersion: v1
kind: Namespace
metadata:
  name: kube-monitoring
---
apiVersion: v1
kind: ConfigMap
metadata:
  labels:
    app: prometheus-grafana
    release: prometheus
  name: prometheus-grafana-resources
  namespace: kube-monitoring
data:
  nodes-dashboard.json: |-
    {
        "dashboard": {
            "__inputs": [
                {
                    "description": "",
                    "label": "prometheus",
                    "name": "DS_PROMETHEUS",
                    "pluginId": "prometheus",
                    "pluginName": "Prometheus",
                    "type": "datasource"
                }
            ],
            "annotations": {
                "list": []
            },
            "description": "Dashboard to get an overview of one server",
            "editable": true,
            "gnetId": 22,
            "graphTooltip": 0,
            "hideControls": false,
            "links": [],
            "refresh": false,
            "rows": [
                {
                    "collapse": false,
                    "editable": true,
                    "height": "250px",
                    "panels": [
                        {
                            "aliasColors": {},
                            "bars": false,
                            "dashLength": 10,
                            "dashes": false,
                            "datasource": "${DS_PROMETHEUS}",
                            "editable": true,
                            "error": false,
                            "fill": 1,
                            "grid": {
                                "threshold1Color": "rgba(216, 200, 27, 0.27)",
                                "threshold2Color": "rgba(234, 112, 112, 0.22)"
                            },
                            "id": 3,
                            "isNew": false,
                            "legend": {
                                "alignAsTable": false,
                                "avg": false,
                                "current": false,
                                "hideEmpty": false,
                                "hideZero": false,
                                "max": false,
                                "min": false,
                                "rightSide": false,
                                "show": true,
                                "total": false
                            },
                            "lines": true,
                            "linewidth": 2,
                            "links": [],
                            "nullPointMode": "connected",
                            "percentage": false,
                            "pointradius": 5,
                            "points": false,
                            "renderer": "flot",
                            "seriesOverrides": [],
                            "spaceLength": 10,
                            "span": 6,
                            "stack": false,
                            "steppedLine": false,
                            "targets": [
                                {
                                    "expr": "100 - (avg by (cpu) (irate(node_cpu{mode=\"idle\", instance=\"$server\"}[5m])) * 100)",
                                    "hide": false,
                                    "intervalFactor": 10,
                                    "legendFormat": "{{cpu}}",
                                    "refId": "A",
                                    "step": 50
                                }
                            ],
                            "title": "Idle CPU",
                            "tooltip": {
                                "msResolution": false,
                                "shared": true,
                                "sort": 0,
                                "value_type": "cumulative"
                            },
                            "type": "graph",
                            "xaxis": {
                                "mode": "time",
                                "show": true,
                                "values": []
                            },
                            "yaxes": [
                                {
                                    "format": "percent",
                                    "label": "cpu usage",
                                    "logBase": 1,
                                    "max": 100,
                                    "min": 0,
                                    "show": true
                                },
                                {
                                    "format": "short",
                                    "logBase": 1,
                                    "show": true
                                }
                            ]
                        },
                        {
                            "aliasColors": {},
                            "bars": false,
                            "dashLength": 10,
                            "dashes": false,
                            "datasource": "${DS_PROMETHEUS}",
                            "editable": true,
                            "error": false,
                            "fill": 1,
                            "grid": {
                                "threshold1Color": "rgba(216, 200, 27, 0.27)",
                                "threshold2Color": "rgba(234, 112, 112, 0.22)"
                            },
                            "id": 9,
                            "isNew": false,
                            "legend": {
                                "alignAsTable": false,
                                "avg": false,
                                "current": false,
                                "hideEmpty": false,
                                "hideZero": false,
                                "max": false,
                                "min": false,
                                "rightSide": false,
                                "show": true,
                                "total": false
                            },
                            "lines": true,
                            "linewidth": 2,
                            "links": [],
                            "nullPointMode": "connected",
                            "percentage": false,
                            "pointradius": 5,
                            "points": false,
                            "renderer": "flot",
                            "seriesOverrides": [],
                            "spaceLength": 10,
                            "span": 6,
                            "stack": false,
                            "steppedLine": false,
                            "targets": [
                                {
                                    "expr": "node_load1{instance=\"$server\"}",
                                    "intervalFactor": 4,
                                    "legendFormat": "load 1m",
                                    "refId": "A",
                                    "step": 20,
                                    "target": ""
                                },
                                {
                                    "expr": "node_load5{instance=\"$server\"}",
                                    "intervalFactor": 4,
                                    "legendFormat": "load 5m",
                                    "refId": "B",
                                    "step": 20,
                                    "target": ""
                                },
                                {
                                    "expr": "node_load15{instance=\"$server\"}",
                                    "intervalFactor": 4,
                                    "legendFormat": "load 15m",
                                    "refId": "C",
                                    "step": 20,
                                    "target": ""
                                }
                            ],
                            "title": "System Load",
                            "tooltip": {
                                "msResolution": false,
                                "shared": true,
                                "sort": 0,
                                "value_type": "cumulative"
                            },
                            "type": "graph",
                            "xaxis": {
                                "mode": "time",
                                "show": true,
                                "values": []
                            },
                            "yaxes": [
                                {
                                    "format": "percentunit",
                                    "logBase": 1,
                                    "show": true
                                },
                                {
                                    "format": "short",
                                    "logBase": 1,
                                    "show": true
                                }
                            ]
                        }
                    ],
                    "showTitle": false,
                    "title": "New Row",
                    "titleSize": "h6"
                },
                {
                    "collapse": false,
                    "editable": true,
                    "height": "250px",
                    "panels": [
                        {
                            "aliasColors": {},
                            "bars": false,
                            "dashLength": 10,
                            "dashes": false,
                            "datasource": "${DS_PROMETHEUS}",
                            "editable": true,
                            "error": false,
                            "fill": 1,
                            "grid": {
                                "threshold1Color": "rgba(216, 200, 27, 0.27)",
                                "threshold2Color": "rgba(234, 112, 112, 0.22)"
                            },
                            "id": 4,
                            "isNew": false,
                            "legend": {
                                "alignAsTable": false,
                                "avg": false,
                                "current": false,
                                "hideEmpty": false,
                                "hideZero": false,
                                "max": false,
                                "min": false,
                                "rightSide": false,
                                "show": true,
                                "total": false
                            },
                            "lines": true,
                            "linewidth": 2,
                            "links": [],
                            "nullPointMode": "connected",
                            "percentage": false,
                            "pointradius": 5,
                            "points": false,
                            "renderer": "flot",
                            "seriesOverrides": [
                                {
                                    "alias": "node_memory_SwapFree{instance=\"172.17.0.1:9100\",job=\"prometheus\"}",
                                    "yaxis": 2
                                }
                            ],
                            "spaceLength": 10,
                            "span": 9,
                            "stack": true,
                            "steppedLine": false,
                            "targets": [
                                {
                                    "expr": "node_memory_MemTotal{instance=\"$server\"} - node_memory_MemFree{instance=\"$server\"} - node_memory_Buffers{instance=\"$server\"} - node_memory_Cached{instance=\"$server\"}",
                                    "hide": false,
                                    "interval": "",
                                    "intervalFactor": 2,
                                    "legendFormat": "memory used",
                                    "metric": "",
                                    "refId": "C",
                                    "step": 10
                                },
                                {
                                    "expr": "node_memory_Buffers{instance=\"$server\"}",
                                    "interval": "",
                                    "intervalFactor": 2,
                                    "legendFormat": "memory buffers",
                                    "metric": "",
                                    "refId": "E",
                                    "step": 10
                                },
                                {
                                    "expr": "node_memory_Cached{instance=\"$server\"}",
                                    "intervalFactor": 2,
                                    "legendFormat": "memory cached",
                                    "metric": "",
                                    "refId": "F",
                                    "step": 10
                                },
                                {
                                    "expr": "node_memory_MemFree{instance=\"$server\"}",
                                    "intervalFactor": 2,
                                    "legendFormat": "memory free",
                                    "metric": "",
                                    "refId": "D",
                                    "step": 10
                                }
                            ],
                            "title": "Memory Usage",
                            "tooltip": {
                                "msResolution": false,
                                "shared": true,
                                "sort": 0,
                                "value_type": "individual"
                            },
                            "type": "graph",
                            "xaxis": {
                                "mode": "time",
                                "show": true,
                                "values": []
                            },
                            "yaxes": [
                                {
                                    "format": "bytes",
                                    "logBase": 1,
                                    "min": "0",
                                    "show": true
                                },
                                {
                                    "format": "short",
                                    "logBase": 1,
                                    "show": true
                                }
                            ]
                        },
                        {
                            "colorBackground": false,
                            "colorValue": false,
                            "colors": [
                                "rgba(50, 172, 45, 0.97)",
                                "rgba(237, 129, 40, 0.89)",
                                "rgba(245, 54, 54, 0.9)"
                            ],
                            "datasource": "${DS_PROMETHEUS}",
                            "editable": true,
                            "format": "percent",
                            "gauge": {
                                "maxValue": 100,
                                "minValue": 0,
                                "show": true,
                                "thresholdLabels": false,
                                "thresholdMarkers": true
                            },
                            "hideTimeOverride": false,
                            "id": 5,
                            "links": [],
                            "mappingType": 1,
                            "mappingTypes": [
                                {
                                    "name": "value to text",
                                    "value": 1
                                },
                                {
                                    "name": "range to text",
                                    "value": 2
                                }
                            ],
                            "maxDataPoints": 100,
                            "nullPointMode": "connected",
                            "postfix": "",
                            "postfixFontSize": "50%",
                            "prefix": "",
                            "prefixFontSize": "50%",
                            "rangeMaps": [
                                {
                                    "from": "null",
                                    "text": "N/A",
                                    "to": "null"
                                }
                            ],
                            "span": 3,
                            "sparkline": {
                                "fillColor": "rgba(31, 118, 189, 0.18)",
                                "full": false,
                                "lineColor": "rgb(31, 120, 193)",
                                "show": false
                            },
                            "targets": [
                                {
                                    "expr": "((node_memory_MemTotal{instance=\"$server\"} - node_memory_MemFree{instance=\"$server\"}  - node_memory_Buffers{instance=\"$server\"} - node_memory_Cached{instance=\"$server\"}) / node_memory_MemTotal{instance=\"$server\"}) * 100",
                                    "intervalFactor": 2,
                                    "refId": "A",
                                    "step": 60,
                                    "target": ""
                                }
                            ],
                            "thresholds": "80, 90",
                            "title": "Memory Usage",
                            "transparent": false,
                            "type": "singlestat",
                            "valueFontSize": "80%",
                            "valueMaps": [
                                {
                                    "op": "=",
                                    "text": "N/A",
                                    "value": "null"
                                }
                            ],
                            "valueName": "avg"
                        }
                    ],
                    "showTitle": false,
                    "title": "New Row",
                    "titleSize": "h6"
                },
                {
                    "collapse": false,
                    "editable": true,
                    "height": "250px",
                    "panels": [
                        {
                            "aliasColors": {},
                            "bars": false,
                            "dashLength": 10,
                            "dashes": false,
                            "datasource": "${DS_PROMETHEUS}",
                            "editable": true,
                            "error": false,
                            "fill": 1,
                            "grid": {
                                "threshold1Color": "rgba(216, 200, 27, 0.27)",
                                "threshold2Color": "rgba(234, 112, 112, 0.22)"
                            },
                            "id": 6,
                            "isNew": true,
                            "legend": {
                                "alignAsTable": false,
                                "avg": false,
                                "current": false,
                                "hideEmpty": false,
                                "hideZero": false,
                                "max": false,
                                "min": false,
                                "rightSide": false,
                                "show": true,
                                "total": false
                            },
                            "lines": true,
                            "linewidth": 2,
                            "links": [],
                            "nullPointMode": "connected",
                            "percentage": false,
                            "pointradius": 5,
                            "points": false,
                            "renderer": "flot",
                            "seriesOverrides": [
                                {
                                    "alias": "read",
                                    "yaxis": 1
                                },
                                {
                                    "alias": "{instance=\"172.17.0.1:9100\"}",
                                    "yaxis": 2
                                },
                                {
                                    "alias": "io time",
                                    "yaxis": 2
                                }
                            ],
                            "spaceLength": 10,
                            "span": 9,
                            "stack": false,
                            "steppedLine": false,
                            "targets": [
                                {
                                    "expr": "sum by (instance) (rate(node_disk_bytes_read{instance=\"$server\"}[2m]))",
                                    "hide": false,
                                    "intervalFactor": 4,
                                    "legendFormat": "read",
                                    "refId": "A",
                                    "step": 20,
                                    "target": ""
                                },
                                {
                                    "expr": "sum by (instance) (rate(node_disk_bytes_written{instance=\"$server\"}[2m]))",
                                    "intervalFactor": 4,
                                    "legendFormat": "written",
                                    "refId": "B",
                                    "step": 20
                                },
                                {
                                    "expr": "sum by (instance) (rate(node_disk_io_time_ms{instance=\"$server\"}[2m]))",
                                    "intervalFactor": 4,
                                    "legendFormat": "io time",
                                    "refId": "C",
                                    "step": 20
                                }
                            ],
                            "title": "Disk I/O",
                            "tooltip": {
                                "msResolution": false,
                                "shared": true,
                                "sort": 0,
                                "value_type": "cumulative"
                            },
                            "type": "graph",
                            "xaxis": {
                                "mode": "time",
                                "show": true,
                                "values": []
                            },
                            "yaxes": [
                                {
                                    "format": "bytes",
                                    "logBase": 1,
                                    "show": true
                                },
                                {
                                    "format": "ms",
                                    "logBase": 1,
                                    "show": true
                                }
                            ]
                        },
                        {
                            "colorBackground": false,
                            "colorValue": false,
                            "colors": [
                                "rgba(50, 172, 45, 0.97)",
                                "rgba(237, 129, 40, 0.89)",
                                "rgba(245, 54, 54, 0.9)"
                            ],
                            "datasource": "${DS_PROMETHEUS}",
                            "editable": true,
                            "format": "percentunit",
                            "gauge": {
                                "maxValue": 1,
                                "minValue": 0,
                                "show": true,
                                "thresholdLabels": false,
                                "thresholdMarkers": true
                            },
                            "hideTimeOverride": false,
                            "id": 7,
                            "links": [],
                            "mappingType": 1,
                            "mappingTypes": [
                                {
                                    "name": "value to text",
                                    "value": 1
                                },
                                {
                                    "name": "range to text",
                                    "value": 2
                                }
                            ],
                            "maxDataPoints": 100,
                            "nullPointMode": "connected",
                            "postfix": "",
                            "postfixFontSize": "50%",
                            "prefix": "",
                            "prefixFontSize": "50%",
                            "rangeMaps": [
                                {
                                    "from": "null",
                                    "text": "N/A",
                                    "to": "null"
                                }
                            ],
                            "span": 3,
                            "sparkline": {
                                "fillColor": "rgba(31, 118, 189, 0.18)",
                                "full": false,
                                "lineColor": "rgb(31, 120, 193)",
                                "show": false
                            },
                            "targets": [
                                {
                                    "expr": "(sum(node_filesystem_size{device!=\"rootfs\",instance=\"$server\"}) - sum(node_filesystem_free{device!=\"rootfs\",instance=\"$server\"})) / sum(node_filesystem_size{device!=\"rootfs\",instance=\"$server\"})",
                                    "intervalFactor": 2,
                                    "refId": "A",
                                    "step": 60,
                                    "target": ""
                                }
                            ],
                            "thresholds": "0.75, 0.9",
                            "title": "Disk Space Usage",
                            "transparent": false,
                            "type": "singlestat",
                            "valueFontSize": "80%",
                            "valueMaps": [
                                {
                                    "op": "=",
                                    "text": "N/A",
                                    "value": "null"
                                }
                            ],
                            "valueName": "current"
                        }
                    ],
                    "showTitle": false,
                    "title": "New Row",
                    "titleSize": "h6"
                },
                {
                    "collapse": false,
                    "editable": true,
                    "height": "250px",
                    "panels": [
                        {
                            "aliasColors": {},
                            "bars": false,
                            "dashLength": 10,
                            "dashes": false,
                            "datasource": "${DS_PROMETHEUS}",
                            "editable": true,
                            "error": false,
                            "fill": 1,
                            "grid": {
                                "threshold1Color": "rgba(216, 200, 27, 0.27)",
                                "threshold2Color": "rgba(234, 112, 112, 0.22)"
                            },
                            "id": 8,
                            "isNew": false,
                            "legend": {
                                "alignAsTable": false,
                                "avg": false,
                                "current": false,
                                "hideEmpty": false,
                                "hideZero": false,
                                "max": false,
                                "min": false,
                                "rightSide": false,
                                "show": true,
                                "total": false
                            },
                            "lines": true,
                            "linewidth": 2,
                            "links": [],
                            "nullPointMode": "connected",
                            "percentage": false,
                            "pointradius": 5,
                            "points": false,
                            "renderer": "flot",
                            "seriesOverrides": [
                                {
                                    "alias": "transmitted",
                                    "yaxis": 2
                                }
                            ],
                            "spaceLength": 10,
                            "span": 6,
                            "stack": false,
                            "steppedLine": false,
                            "targets": [
                                {
                                    "expr": "rate(node_network_receive_bytes{instance=\"$server\",device!~\"lo\"}[5m])",
                                    "hide": false,
                                    "intervalFactor": 2,
                                    "legendFormat": "{{device}}",
                                    "refId": "A",
                                    "step": 10,
                                    "target": ""
                                }
                            ],
                            "title": "Network Received",
                            "tooltip": {
                                "msResolution": false,
                                "shared": true,
                                "sort": 0,
                                "value_type": "cumulative"
                            },
                            "type": "graph",
                            "xaxis": {
                                "mode": "time",
                                "show": true,
                                "values": []
                            },
                            "yaxes": [
                                {
                                    "format": "bytes",
                                    "logBase": 1,
                                    "show": true
                                },
                                {
                                    "format": "bytes",
                                    "logBase": 1,
                                    "show": true
                                }
                            ]
                        },
                        {
                            "aliasColors": {},
                            "bars": false,
                            "dashLength": 10,
                            "dashes": false,
                            "datasource": "${DS_PROMETHEUS}",
                            "editable": true,
                            "error": false,
                            "fill": 1,
                            "grid": {
                                "threshold1Color": "rgba(216, 200, 27, 0.27)",
                                "threshold2Color": "rgba(234, 112, 112, 0.22)"
                            },
                            "id": 10,
                            "isNew": false,
                            "legend": {
                                "alignAsTable": false,
                                "avg": false,
                                "current": false,
                                "hideEmpty": false,
                                "hideZero": false,
                                "max": false,
                                "min": false,
                                "rightSide": false,
                                "show": true,
                                "total": false
                            },
                            "lines": true,
                            "linewidth": 2,
                            "links": [],
                            "nullPointMode": "connected",
                            "percentage": false,
                            "pointradius": 5,
                            "points": false,
                            "renderer": "flot",
                            "seriesOverrides": [
                                {
                                    "alias": "transmitted",
                                    "yaxis": 2
                                }
                            ],
                            "spaceLength": 10,
                            "span": 6,
                            "stack": false,
                            "steppedLine": false,
                            "targets": [
                                {
                                    "expr": "rate(node_network_transmit_bytes{instance=\"$server\",device!~\"lo\"}[5m])",
                                    "hide": false,
                                    "intervalFactor": 2,
                                    "legendFormat": "{{device}}",
                                    "refId": "B",
                                    "step": 10,
                                    "target": ""
                                }
                            ],
                            "title": "Network Transmitted",
                            "tooltip": {
                                "msResolution": false,
                                "shared": true,
                                "sort": 0,
                                "value_type": "cumulative"
                            },
                            "type": "graph",
                            "xaxis": {
                                "mode": "time",
                                "show": true,
                                "values": []
                            },
                            "yaxes": [
                                {
                                    "format": "bytes",
                                    "logBase": 1,
                                    "show": true
                                },
                                {
                                    "format": "bytes",
                                    "logBase": 1,
                                    "show": true
                                }
                            ]
                        }
                    ],
                    "showTitle": false,
                    "title": "New Row",
                    "titleSize": "h6"
                }
            ],
            "schemaVersion": 14,
            "sharedCrosshair": false,
            "style": "dark",
            "tags": [],
            "templating": {
                "list": [
                    {
                        "allValue": null,
                        "current": {},
                        "datasource": "${DS_PROMETHEUS}",
                        "hide": 0,
                        "includeAll": false,
                        "label": null,
                        "multi": false,
                        "name": "server",
                        "options": [],
                        "query": "label_values(node_boot_time, instance)",
                        "refresh": 1,
                        "regex": "",
                        "sort": 0,
                        "tagValuesQuery": "",
                        "tags": [],
                        "tagsQuery": "",
                        "type": "query",
                        "useTags": false
                    }
                ]
            },
            "time": {
                "from": "now-1h",
                "to": "now"
            },
            "timepicker": {
                "refresh_intervals": [
                    "5s",
                    "10s",
                    "30s",
                    "1m",
                    "5m",
                    "15m",
                    "30m",
                    "1h",
                    "2h",
                    "1d"
                ],
                "time_options": [
                    "5m",
                    "15m",
                    "1h",
                    "6h",
                    "12h",
                    "24h",
                    "2d",
                    "7d",
                    "30d"
                ]
            },
            "timezone": "browser",
            "title": "Kubernetes Nodes",
            "version": 2
        },
        "inputs": [
            {
                "name": "DS_PROMETHEUS",
                "pluginId": "prometheus",
                "type": "datasource",
                "value": "prometheus"
            }
        ],
        "overwrite": true
    }
  pods-dashboard.json: |
    {
        "dashboard": {
            "__inputs": [
                {
                    "description": "",
                    "label": "prometheus",
                    "name": "DS_PROMETHEUS",
                    "pluginId": "prometheus",
                    "pluginName": "Prometheus",
                    "type": "datasource"
                }
            ],
            "annotations": {
                "list": []
            },
            "editable": true,
            "graphTooltip": 1,
            "hideControls": false,
            "links": [],
            "refresh": false,
            "rows": [
                {
                    "collapse": false,
                    "editable": true,
                    "height": "250px",
                    "panels": [
                        {
                            "aliasColors": {},
                            "bars": false,
                            "dashLength": 10,
                            "dashes": false,
                            "datasource": "${DS_PROMETHEUS}",
                            "editable": true,
                            "error": false,
                            "fill": 1,
                            "grid": {
                                "threshold1Color": "rgba(216, 200, 27, 0.27)",
                                "threshold2Color": "rgba(234, 112, 112, 0.22)"
                            },
                            "id": 1,
                            "isNew": false,
                            "legend": {
                                "alignAsTable": true,
                                "avg": true,
                                "current": true,
                                "hideEmpty": false,
                                "hideZero": false,
                                "max": false,
                                "min": false,
                                "rightSide": true,
                                "show": true,
                                "total": false,
                                "values": true
                            },
                            "lines": true,
                            "linewidth": 2,
                            "links": [],
                            "nullPointMode": "connected",
                            "percentage": false,
                            "pointradius": 5,
                            "points": false,
                            "renderer": "flot",
                            "seriesOverrides": [],
                            "spaceLength": 10,
                            "span": 12,
                            "stack": false,
                            "steppedLine": false,
                            "targets": [
                                {
                                    "expr": "sum by(container_name) (container_memory_usage_bytes{pod_name=\"$pod\", container_name=~\"$container\", container_name!=\"POD\"})",
                                    "interval": "10s",
                                    "intervalFactor": 1,
                                    "legendFormat": "Current: {{ container_name }}",
                                    "metric": "container_memory_usage_bytes",
                                    "refId": "A",
                                    "step": 15
                                },
                                {
                                    "expr": "kube_pod_container_resource_requests_memory_bytes{pod=\"$pod\", container=~\"$container\"}",
                                    "interval": "10s",
                                    "intervalFactor": 2,
                                    "legendFormat": "Requested: {{ container }}",
                                    "metric": "kube_pod_container_resource_requests_memory_bytes",
                                    "refId": "B",
                                    "step": 20
                                },
                                {
                                    "expr": "kube_pod_container_resource_limits_memory_bytes{pod=\"$pod\", container=~\"$container\"}",
                                    "interval": "10s",
                                    "intervalFactor": 2,
                                    "legendFormat": "Limit: {{ container }}",
                                    "metric": "kube_pod_container_resource_limits_memory_bytes",
                                    "refId": "C",
                                    "step": 20
                                }
                            ],
                            "title": "Memory Usage",
                            "tooltip": {
                                "msResolution": true,
                                "shared": true,
                                "sort": 0,
                                "value_type": "cumulative"
                            },
                            "type": "graph",
                            "xaxis": {
                                "mode": "time",
                                "show": true,
                                "values": []
                            },
                            "yaxes": [
                                {
                                    "format": "bytes",
                                    "logBase": 1,
                                    "show": true
                                },
                                {
                                    "format": "short",
                                    "logBase": 1,
                                    "show": true
                                }
                            ]
                        }
                    ],
                    "showTitle": false,
                    "title": "Row",
                    "titleSize": "h6"
                },
                {
                    "collapse": false,
                    "editable": true,
                    "height": "250px",
                    "panels": [
                        {
                            "aliasColors": {},
                            "bars": false,
                            "dashLength": 10,
                            "dashes": false,
                            "datasource": "${DS_PROMETHEUS}",
                            "editable": true,
                            "error": false,
                            "fill": 1,
                            "grid": {
                                "threshold1Color": "rgba(216, 200, 27, 0.27)",
                                "threshold2Color": "rgba(234, 112, 112, 0.22)"
                            },
                            "id": 2,
                            "isNew": false,
                            "legend": {
                                "alignAsTable": true,
                                "avg": true,
                                "current": true,
                                "hideEmpty": false,
                                "hideZero": false,
                                "max": false,
                                "min": false,
                                "rightSide": true,
                                "show": true,
                                "total": false,
                                "values": true
                            },
                            "lines": true,
                            "linewidth": 2,
                            "links": [],
                            "nullPointMode": "connected",
                            "percentage": false,
                            "pointradius": 5,
                            "points": false,
                            "renderer": "flot",
                            "seriesOverrides": [],
                            "spaceLength": 10,
                            "span": 12,
                            "stack": false,
                            "steppedLine": false,
                            "targets": [
                                {
                                    "expr": "sum by (container_name)(rate(container_cpu_usage_seconds_total{image!=\"\",container_name!=\"POD\",pod_name=\"$pod\"}[5m]))",
                                    "intervalFactor": 2,
                                    "legendFormat": "{{ container_name }}",
                                    "refId": "A",
                                    "step": 30
                                },
                                {
                                    "expr": "kube_pod_container_resource_requests_cpu_cores{pod=\"$pod\", container=~\"$container\"}",
                                    "interval": "10s",
                                    "intervalFactor": 2,
                                    "legendFormat": "Requested: {{ container }}",
                                    "metric": "kube_pod_container_resource_requests_cpu_cores",
                                    "refId": "B",
                                    "step": 20
                                },
                                {
                                    "expr": "kube_pod_container_resource_limits_cpu_cores{pod=\"$pod\", container=~\"$container\"}",
                                    "interval": "10s",
                                    "intervalFactor": 2,
                                    "legendFormat": "Limit: {{ container }}",
                                    "metric": "kube_pod_container_resource_limits_memory_bytes",
                                    "refId": "C",
                                    "step": 20
                                }
                            ],
                            "title": "CPU Usage",
                            "tooltip": {
                                "msResolution": true,
                                "shared": true,
                                "sort": 0,
                                "value_type": "cumulative"
                            },
                            "type": "graph",
                            "xaxis": {
                                "mode": "time",
                                "show": true,
                                "values": []
                            },
                            "yaxes": [
                                {
                                    "format": "short",
                                    "logBase": 1,
                                    "show": true
                                },
                                {
                                    "format": "short",
                                    "logBase": 1,
                                    "show": true
                                }
                            ]
                        }
                    ],
                    "showTitle": false,
                    "title": "Row",
                    "titleSize": "h6"
                },
                {
                    "collapse": false,
                    "editable": true,
                    "height": "250px",
                    "panels": [
                        {
                            "aliasColors": {},
                            "bars": false,
                            "dashLength": 10,
                            "dashes": false,
                            "datasource": "${DS_PROMETHEUS}",
                            "editable": true,
                            "error": false,
                            "fill": 1,
                            "grid": {
                                "threshold1Color": "rgba(216, 200, 27, 0.27)",
                                "threshold2Color": "rgba(234, 112, 112, 0.22)"
                            },
                            "id": 3,
                            "isNew": false,
                            "legend": {
                                "alignAsTable": true,
                                "avg": true,
                                "current": true,
                                "hideEmpty": false,
                                "hideZero": false,
                                "max": false,
                                "min": false,
                                "rightSide": true,
                                "show": true,
                                "total": false,
                                "values": true
                            },
                            "lines": true,
                            "linewidth": 2,
                            "links": [],
                            "nullPointMode": "connected",
                            "percentage": false,
                            "pointradius": 5,
                            "points": false,
                            "renderer": "flot",
                            "seriesOverrides": [],
                            "spaceLength": 10,
                            "span": 12,
                            "stack": false,
                            "steppedLine": false,
                            "targets": [
                                {
                                    "expr": "sort_desc(sum by (pod_name) (rate(container_network_receive_bytes_total{pod_name=\"$pod\"}[5m])))",
                                    "intervalFactor": 2,
                                    "legendFormat": "{{ pod_name }}",
                                    "refId": "A",
                                    "step": 30
                                }
                            ],
                            "title": "Network I/O",
                            "tooltip": {
                                "msResolution": true,
                                "shared": true,
                                "sort": 0,
                                "value_type": "cumulative"
                            },
                            "type": "graph",
                            "xaxis": {
                                "mode": "time",
                                "show": true,
                                "values": []
                            },
                            "yaxes": [
                                {
                                    "format": "bytes",
                                    "logBase": 1,
                                    "show": true
                                },
                                {
                                    "format": "short",
                                    "logBase": 1,
                                    "show": true
                                }
                            ]
                        }
                    ],
                    "showTitle": false,
                    "title": "New Row",
                    "titleSize": "h6"
                }
            ],
            "schemaVersion": 14,
            "sharedCrosshair": false,
            "style": "dark",
            "tags": [],
            "templating": {
                "list": [
                    {
                        "allValue": ".*",
                        "current": {},
                        "datasource": "${DS_PROMETHEUS}",
                        "hide": 0,
                        "includeAll": true,
                        "label": "Namespace",
                        "multi": false,
                        "name": "namespace",
                        "options": [],
                        "query": "label_values(kube_pod_info, namespace)",
                        "refresh": 1,
                        "regex": "",
                        "sort": 0,
                        "tagValuesQuery": "",
                        "tags": [],
                        "tagsQuery": "",
                        "type": "query",
                        "useTags": false
                    },
                    {
                        "allValue": null,
                        "current": {},
                        "datasource": "${DS_PROMETHEUS}",
                        "hide": 0,
                        "includeAll": false,
                        "label": "Pod",
                        "multi": false,
                        "name": "pod",
                        "options": [],
                        "query": "label_values(kube_pod_info{namespace=~\"$namespace\"}, pod)",
                        "refresh": 1,
                        "regex": "",
                        "sort": 0,
                        "tagValuesQuery": "",
                        "tags": [],
                        "tagsQuery": "",
                        "type": "query",
                        "useTags": false
                    },
                    {
                        "allValue": ".*",
                        "current": {},
                        "datasource": "${DS_PROMETHEUS}",
                        "hide": 0,
                        "includeAll": true,
                        "label": "Container",
                        "multi": false,
                        "name": "container",
                        "options": [],
                        "query": "label_values(kube_pod_container_info{namespace=\"$namespace\", pod=\"$pod\"}, container)",
                        "refresh": 1,
                        "regex": "",
                        "sort": 0,
                        "tagValuesQuery": "",
                        "tags": [],
                        "tagsQuery": "",
                        "type": "query",
                        "useTags": false
                    }
                ]
            },
            "time": {
                "from": "now-6h",
                "to": "now"
            },
            "timepicker": {
                "refresh_intervals": [
                    "5s",
                    "10s",
                    "30s",
                    "1m",
                    "5m",
                    "15m",
                    "30m",
                    "1h",
                    "2h",
                    "1d"
                ],
                "time_options": [
                    "5m",
                    "15m",
                    "1h",
                    "6h",
                    "12h",
                    "24h",
                    "2d",
                    "7d",
                    "30d"
                ]
            },
            "timezone": "browser",
            "title": "Kubernetes Pods",
            "version": 1
        },
        "inputs": [
            {
                "name": "DS_PROMETHEUS",
                "pluginId": "prometheus",
                "type": "datasource",
                "value": "prometheus"
            }
        ],
        "overwrite": true
    }
---
apiVersion: v1
kind: ConfigMap
metadata:
  labels:
    app: prometheus-grafana
    release: prometheus
  name: prometheus-grafana
  namespace: kube-monitoring
data:
  kubernetes-capacity-planning-dashboard.json: |-
    {
        "dashboard": {
            "__inputs": [
                {
                    "description": "",
                    "label": "prometheus",
                    "name": "DS_PROMETHEUS",
                    "pluginId": "prometheus",
                    "pluginName": "Prometheus",
                    "type": "datasource"
                }
            ],
            "annotations": {
                "list": []
            },
            "editable": true,
            "gnetId": 22,
            "graphTooltip": 0,
            "hideControls": false,
            "links": [],
            "refresh": false,
            "rows": [
                {
                    "collapse": false,
                    "editable": true,
                    "height": "250px",
                    "panels": [
                        {
                            "aliasColors": {},
                            "bars": false,
                            "dashLength": 10,
                            "dashes": false,
                            "datasource": "${DS_PROMETHEUS}",
                            "editable": true,
                            "error": false,
                            "fill": 1,
                            "grid": {
                                "threshold1Color": "rgba(216, 200, 27, 0.27)",
                                "threshold2Color": "rgba(234, 112, 112, 0.22)"
                            },
                            "id": 3,
                            "isNew": false,
                            "legend": {
                                "alignAsTable": false,
                                "avg": false,
                                "current": false,
                                "hideEmpty": false,
                                "hideZero": false,
                                "max": false,
                                "min": false,
                                "rightSide": false,
                                "show": true,
                                "total": false
                            },
                            "lines": true,
                            "linewidth": 2,
                            "links": [],
                            "nullPointMode": "connected",
                            "percentage": false,
                            "pointradius": 5,
                            "points": false,
                            "renderer": "flot",
                            "seriesOverrides": [],
                            "spaceLength": 10,
                            "span": 6,
                            "stack": false,
                            "steppedLine": false,
                            "targets": [
                                {
                                    "expr": "sum(rate(node_cpu{mode=\"idle\"}[2m])) * 100",
                                    "hide": false,
                                    "intervalFactor": 10,
                                    "legendFormat": "",
                                    "refId": "A",
                                    "step": 50
                                }
                            ],
                            "title": "Idle CPU",
                            "tooltip": {
                                "msResolution": false,
                                "shared": true,
                                "sort": 0,
                                "value_type": "cumulative"
                            },
                            "type": "graph",
                            "xaxis": {
                                "mode": "time",
                                "show": true,
                                "values": []
                            },
                            "yaxes": [
                                {
                                    "format": "percent",
                                    "label": "cpu usage",
                                    "logBase": 1,
                                    "min": 0,
                                    "show": true
                                },
                                {
                                    "format": "short",
                                    "logBase": 1,
                                    "show": true
                                }
                            ]
                        },
                        {
                            "aliasColors": {},
                            "bars": false,
                            "dashLength": 10,
                            "dashes": false,
                            "datasource": "${DS_PROMETHEUS}",
                            "editable": true,
                            "error": false,
                            "fill": 1,
                            "grid": {
                                "threshold1Color": "rgba(216, 200, 27, 0.27)",
                                "threshold2Color": "rgba(234, 112, 112, 0.22)"
                            },
                            "id": 9,
                            "isNew": false,
                            "legend": {
                                "alignAsTable": false,
                                "avg": false,
                                "current": false,
                                "hideEmpty": false,
                                "hideZero": false,
                                "max": false,
                                "min": false,
                                "rightSide": false,
                                "show": true,
                                "total": false
                            },
                            "lines": true,
                            "linewidth": 2,
                            "links": [],
                            "nullPointMode": "connected",
                            "percentage": false,
                            "pointradius": 5,
                            "points": false,
                            "renderer": "flot",
                            "seriesOverrides": [],
                            "spaceLength": 10,
                            "span": 6,
                            "stack": false,
                            "steppedLine": false,
                            "targets": [
                                {
                                    "expr": "sum(node_load1)",
                                    "intervalFactor": 4,
                                    "legendFormat": "load 1m",
                                    "refId": "A",
                                    "step": 20,
                                    "target": ""
                                },
                                {
                                    "expr": "sum(node_load5)",
                                    "intervalFactor": 4,
                                    "legendFormat": "load 5m",
                                    "refId": "B",
                                    "step": 20,
                                    "target": ""
                                },
                                {
                                    "expr": "sum(node_load15)",
                                    "intervalFactor": 4,
                                    "legendFormat": "load 15m",
                                    "refId": "C",
                                    "step": 20,
                                    "target": ""
                                }
                            ],
                            "title": "System Load",
                            "tooltip": {
                                "msResolution": false,
                                "shared": true,
                                "sort": 0,
                                "value_type": "cumulative"
                            },
                            "type": "graph",
                            "xaxis": {
                                "mode": "time",
                                "show": true,
                                "values": []
                            },
                            "yaxes": [
                                {
                                    "format": "percentunit",
                                    "logBase": 1,
                                    "show": true
                                },
                                {
                                    "format": "short",
                                    "logBase": 1,
                                    "show": true
                                }
                            ]
                        }
                    ],
                    "showTitle": false,
                    "title": "New Row",
                    "titleSize": "h6"
                },
                {
                    "collapse": false,
                    "editable": true,
                    "height": "250px",
                    "panels": [
                        {
                            "aliasColors": {},
                            "bars": false,
                            "dashLength": 10,
                            "dashes": false,
                            "datasource": "${DS_PROMETHEUS}",
                            "editable": true,
                            "error": false,
                            "fill": 1,
                            "grid": {
                                "threshold1Color": "rgba(216, 200, 27, 0.27)",
                                "threshold2Color": "rgba(234, 112, 112, 0.22)"
                            },
                            "id": 4,
                            "isNew": false,
                            "legend": {
                                "alignAsTable": false,
                                "avg": false,
                                "current": false,
                                "hideEmpty": false,
                                "hideZero": false,
                                "max": false,
                                "min": false,
                                "rightSide": false,
                                "show": true,
                                "total": false
                            },
                            "lines": true,
                            "linewidth": 2,
                            "links": [],
                            "nullPointMode": "connected",
                            "percentage": false,
                            "pointradius": 5,
                            "points": false,
                            "renderer": "flot",
                            "seriesOverrides": [
                                {
                                    "alias": "node_memory_SwapFree{instance=\"172.17.0.1:9100\",job=\"prometheus\"}",
                                    "yaxis": 2
                                }
                            ],
                            "spaceLength": 10,
                            "span": 9,
                            "stack": true,
                            "steppedLine": false,
                            "targets": [
                                {
                                    "expr": "sum(node_memory_MemTotal) - sum(node_memory_MemFree) - sum(node_memory_Buffers) - sum(node_memory_Cached)",
                                    "intervalFactor": 2,
                                    "legendFormat": "memory usage",
                                    "metric": "memo",
                                    "refId": "A",
                                    "step": 10,
                                    "target": ""
                                },
                                {
                                    "expr": "sum(node_memory_Buffers)",
                                    "interval": "",
                                    "intervalFactor": 2,
                                    "legendFormat": "memory buffers",
                                    "metric": "memo",
                                    "refId": "B",
                                    "step": 10,
                                    "target": ""
                                },
                                {
                                    "expr": "sum(node_memory_Cached)",
                                    "interval": "",
                                    "intervalFactor": 2,
                                    "legendFormat": "memory cached",
                                    "metric": "memo",
                                    "refId": "C",
                                    "step": 10,
                                    "target": ""
                                },
                                {
                                    "expr": "sum(node_memory_MemFree)",
                                    "interval": "",
                                    "intervalFactor": 2,
                                    "legendFormat": "memory free",
                                    "metric": "memo",
                                    "refId": "D",
                                    "step": 10,
                                    "target": ""
                                }
                            ],
                            "title": "Memory Usage",
                            "tooltip": {
                                "msResolution": false,
                                "shared": true,
                                "sort": 0,
                                "value_type": "individual"
                            },
                            "type": "graph",
                            "xaxis": {
                                "mode": "time",
                                "show": true,
                                "values": []
                            },
                            "yaxes": [
                                {
                                    "format": "bytes",
                                    "logBase": 1,
                                    "min": "0",
                                    "show": true
                                },
                                {
                                    "format": "short",
                                    "logBase": 1,
                                    "show": true
                                }
                            ]
                        },
                        {
                            "colorBackground": false,
                            "colorValue": false,
                            "colors": [
                                "rgba(50, 172, 45, 0.97)",
                                "rgba(237, 129, 40, 0.89)",
                                "rgba(245, 54, 54, 0.9)"
                            ],
                            "datasource": "${DS_PROMETHEUS}",
                            "editable": true,
                            "format": "percent",
                            "gauge": {
                                "maxValue": 100,
                                "minValue": 0,
                                "show": true,
                                "thresholdLabels": false,
                                "thresholdMarkers": true
                            },
                            "hideTimeOverride": false,
                            "id": 5,
                            "links": [],
                            "mappingType": 1,
                            "mappingTypes": [
                                {
                                    "name": "value to text",
                                    "value": 1
                                },
                                {
                                    "name": "range to text",
                                    "value": 2
                                }
                            ],
                            "maxDataPoints": 100,
                            "nullPointMode": "connected",
                            "postfix": "",
                            "postfixFontSize": "50%",
                            "prefix": "",
                            "prefixFontSize": "50%",
                            "rangeMaps": [
                                {
                                    "from": "null",
                                    "text": "N/A",
                                    "to": "null"
                                }
                            ],
                            "span": 3,
                            "sparkline": {
                                "fillColor": "rgba(31, 118, 189, 0.18)",
                                "full": false,
                                "lineColor": "rgb(31, 120, 193)",
                                "show": false
                            },
                            "targets": [
                                {
                                    "expr": "((sum(node_memory_MemTotal) - sum(node_memory_MemFree) - sum(node_memory_Buffers) - sum(node_memory_Cached)) / sum(node_memory_MemTotal)) * 100",
                                    "intervalFactor": 2,
                                    "metric": "",
                                    "refId": "A",
                                    "step": 60,
                                    "target": ""
                                }
                            ],
                            "thresholds": "80, 90",
                            "title": "Memory Usage",
                            "transparent": false,
                            "type": "singlestat",
                            "valueFontSize": "80%",
                            "valueMaps": [
                                {
                                    "op": "=",
                                    "text": "N/A",
                                    "value": "null"
                                }
                            ],
                            "valueName": "avg"
                        }
                    ],
                    "showTitle": false,
                    "title": "New Row",
                    "titleSize": "h6"
                },
                {
                    "collapse": false,
                    "editable": true,
                    "height": "246px",
                    "panels": [
                        {
                            "aliasColors": {},
                            "bars": false,
                            "dashLength": 10,
                            "dashes": false,
                            "datasource": "${DS_PROMETHEUS}",
                            "editable": true,
                            "error": false,
                            "fill": 1,
                            "grid": {
                                "threshold1Color": "rgba(216, 200, 27, 0.27)",
                                "threshold2Color": "rgba(234, 112, 112, 0.22)"
                            },
                            "id": 6,
                            "isNew": false,
                            "legend": {
                                "alignAsTable": false,
                                "avg": false,
                                "current": false,
                                "hideEmpty": false,
                                "hideZero": false,
                                "max": false,
                                "min": false,
                                "rightSide": false,
                                "show": true,
                                "total": false
                            },
                            "lines": true,
                            "linewidth": 2,
                            "links": [],
                            "nullPointMode": "connected",
                            "percentage": false,
                            "pointradius": 5,
                            "points": false,
                            "renderer": "flot",
                            "seriesOverrides": [
                                {
                                    "alias": "read",
                                    "yaxis": 1
                                },
                                {
                                    "alias": "{instance=\"172.17.0.1:9100\"}",
                                    "yaxis": 2
                                },
                                {
                                    "alias": "io time",
                                    "yaxis": 2
                                }
                            ],
                            "spaceLength": 10,
                            "span": 9,
                            "stack": false,
                            "steppedLine": false,
                            "targets": [
                                {
                                    "expr": "sum(rate(node_disk_bytes_read[5m]))",
                                    "hide": false,
                                    "intervalFactor": 4,
                                    "legendFormat": "read",
                                    "refId": "A",
                                    "step": 20,
                                    "target": ""
                                },
                                {
                                    "expr": "sum(rate(node_disk_bytes_written[5m]))",
                                    "intervalFactor": 4,
                                    "legendFormat": "written",
                                    "refId": "B",
                                    "step": 20
                                },
                                {
                                    "expr": "sum(rate(node_disk_io_time_ms[5m]))",
                                    "intervalFactor": 4,
                                    "legendFormat": "io time",
                                    "refId": "C",
                                    "step": 20
                                }
                            ],
                            "title": "Disk I/O",
                            "tooltip": {
                                "msResolution": false,
                                "shared": true,
                                "sort": 0,
                                "value_type": "cumulative"
                            },
                            "type": "graph",
                            "xaxis": {
                                "mode": "time",
                                "show": true,
                                "values": []
                            },
                            "yaxes": [
                                {
                                    "format": "bytes",
                                    "logBase": 1,
                                    "show": true
                                },
                                {
                                    "format": "ms",
                                    "logBase": 1,
                                    "show": true
                                }
                            ]
                        },
                        {
                            "colorBackground": false,
                            "colorValue": false,
                            "colors": [
                                "rgba(50, 172, 45, 0.97)",
                                "rgba(237, 129, 40, 0.89)",
                                "rgba(245, 54, 54, 0.9)"
                            ],
                            "datasource": "${DS_PROMETHEUS}",
                            "editable": true,
                            "format": "percentunit",
                            "gauge": {
                                "maxValue": 1,
                                "minValue": 0,
                                "show": true,
                                "thresholdLabels": false,
                                "thresholdMarkers": true
                            },
                            "hideTimeOverride": false,
                            "id": 12,
                            "links": [],
                            "mappingType": 1,
                            "mappingTypes": [
                                {
                                    "name": "value to text",
                                    "value": 1
                                },
                                {
                                    "name": "range to text",
                                    "value": 2
                                }
                            ],
                            "maxDataPoints": 100,
                            "nullPointMode": "connected",
                            "postfix": "",
                            "postfixFontSize": "50%",
                            "prefix": "",
                            "prefixFontSize": "50%",
                            "rangeMaps": [
                                {
                                    "from": "null",
                                    "text": "N/A",
                                    "to": "null"
                                }
                            ],
                            "span": 3,
                            "sparkline": {
                                "fillColor": "rgba(31, 118, 189, 0.18)",
                                "full": false,
                                "lineColor": "rgb(31, 120, 193)",
                                "show": false
                            },
                            "targets": [
                                {
                                    "expr": "(sum(node_filesystem_size{device!=\"rootfs\"}) - sum(node_filesystem_free{device!=\"rootfs\"})) / sum(node_filesystem_size{device!=\"rootfs\"})",
                                    "intervalFactor": 2,
                                    "refId": "A",
                                    "step": 60,
                                    "target": ""
                                }
                            ],
                            "thresholds": "0.75, 0.9",
                            "title": "Disk Space Usage",
                            "transparent": false,
                            "type": "singlestat",
                            "valueFontSize": "80%",
                            "valueMaps": [
                                {
                                    "op": "=",
                                    "text": "N/A",
                                    "value": "null"
                                }
                            ],
                            "valueName": "current"
                        }
                    ],
                    "showTitle": false,
                    "title": "New Row",
                    "titleSize": "h6"
                },
                {
                    "collapse": false,
                    "editable": true,
                    "height": "250px",
                    "panels": [
                        {
                            "aliasColors": {},
                            "bars": false,
                            "dashLength": 10,
                            "dashes": false,
                            "datasource": "${DS_PROMETHEUS}",
                            "editable": true,
                            "error": false,
                            "fill": 1,
                            "grid": {
                                "threshold1Color": "rgba(216, 200, 27, 0.27)",
                                "threshold2Color": "rgba(234, 112, 112, 0.22)"
                            },
                            "id": 8,
                            "isNew": false,
                            "legend": {
                                "alignAsTable": false,
                                "avg": false,
                                "current": false,
                                "hideEmpty": false,
                                "hideZero": false,
                                "max": false,
                                "min": false,
                                "rightSide": false,
                                "show": true,
                                "total": false
                            },
                            "lines": true,
                            "linewidth": 2,
                            "links": [],
                            "nullPointMode": "connected",
                            "percentage": false,
                            "pointradius": 5,
                            "points": false,
                            "renderer": "flot",
                            "seriesOverrides": [
                                {
                                    "alias": "transmitted",
                                    "yaxis": 2
                                }
                            ],
                            "spaceLength": 10,
                            "span": 6,
                            "stack": false,
                            "steppedLine": false,
                            "targets": [
                                {
                                    "expr": "sum(rate(node_network_receive_bytes{device!~\"lo\"}[5m]))",
                                    "hide": false,
                                    "intervalFactor": 2,
                                    "legendFormat": "",
                                    "refId": "A",
                                    "step": 10,
                                    "target": ""
                                }
                            ],
                            "title": "Network Received",
                            "tooltip": {
                                "msResolution": false,
                                "shared": true,
                                "sort": 0,
                                "value_type": "cumulative"
                            },
                            "type": "graph",
                            "xaxis": {
                                "mode": "time",
                                "show": true,
                                "values": []
                            },
                            "yaxes": [
                                {
                                    "format": "bytes",
                                    "logBase": 1,
                                    "show": true
                                },
                                {
                                    "format": "bytes",
                                    "logBase": 1,
                                    "show": true
                                }
                            ]
                        },
                        {
                            "aliasColors": {},
                            "bars": false,
                            "dashLength": 10,
                            "dashes": false,
                            "datasource": "${DS_PROMETHEUS}",
                            "editable": true,
                            "error": false,
                            "fill": 1,
                            "grid": {
                                "threshold1Color": "rgba(216, 200, 27, 0.27)",
                                "threshold2Color": "rgba(234, 112, 112, 0.22)"
                            },
                            "id": 10,
                            "isNew": false,
                            "legend": {
                                "alignAsTable": false,
                                "avg": false,
                                "current": false,
                                "hideEmpty": false,
                                "hideZero": false,
                                "max": false,
                                "min": false,
                                "rightSide": false,
                                "show": true,
                                "total": false
                            },
                            "lines": true,
                            "linewidth": 2,
                            "links": [],
                            "nullPointMode": "connected",
                            "percentage": false,
                            "pointradius": 5,
                            "points": false,
                            "renderer": "flot",
                            "seriesOverrides": [
                                {
                                    "alias": "transmitted",
                                    "yaxis": 2
                                }
                            ],
                            "spaceLength": 10,
                            "span": 6,
                            "stack": false,
                            "steppedLine": false,
                            "targets": [
                                {
                                    "expr": "sum(rate(node_network_transmit_bytes{device!~\"lo\"}[5m]))",
                                    "hide": false,
                                    "intervalFactor": 2,
                                    "legendFormat": "",
                                    "refId": "B",
                                    "step": 10,
                                    "target": ""
                                }
                            ],
                            "title": "Network Transmitted",
                            "tooltip": {
                                "msResolution": false,
                                "shared": true,
                                "sort": 0,
                                "value_type": "cumulative"
                            },
                            "type": "graph",
                            "xaxis": {
                                "mode": "time",
                                "show": true,
                                "values": []
                            },
                            "yaxes": [
                                {
                                    "format": "bytes",
                                    "logBase": 1,
                                    "show": true
                                },
                                {
                                    "format": "bytes",
                                    "logBase": 1,
                                    "show": true
                                }
                            ]
                        }
                    ],
                    "showTitle": false,
                    "title": "New Row",
                    "titleSize": "h6"
                },
                {
                    "collapse": false,
                    "editable": true,
                    "height": "276px",
                    "panels": [
                        {
                            "aliasColors": {},
                            "bars": false,
                            "dashes": false,
                            "datasource": "${DS_PROMETHEUS}",
                            "editable": true,
                            "error": false,
                            "fill": 1,
                            "grid": {
                                "threshold1Color": "rgba(216, 200, 27, 0.27)",
                                "threshold2Color": "rgba(234, 112, 112, 0.22)"
                            },
                            "id": 11,
                            "isNew": true,
                            "legend": {
                                "alignAsTable": false,
                                "avg": false,
                                "current": false,
                                "hideEmpty": false,
                                "hideZero": false,
                                "max": false,
                                "min": false,
                                "rightSide": false,
                                "show": true,
                                "total": false
                            },
                            "lines": true,
                            "linewidth": 2,
                            "links": [],
                            "nullPointMode": "connected",
                            "percentage": false,
                            "pointradius": 5,
                            "points": false,
                            "renderer": "flot",
                            "seriesOverrides": [],
                            "spaceLength": 11,
                            "span": 9,
                            "stack": false,
                            "steppedLine": false,
                            "targets": [
                                {
                                    "expr": "sum(kube_pod_info)",
                                    "format": "time_series",
                                    "intervalFactor": 2,
                                    "legendFormat": "Current number of Pods",
                                    "refId": "A",
                                    "step": 10
                                },
                                {
                                    "expr": "sum(kube_node_status_capacity_pods)",
                                    "format": "time_series",
                                    "intervalFactor": 2,
                                    "legendFormat": "Maximum capacity of pods",
                                    "refId": "B",
                                    "step": 10
                                }
                            ],
                            "title": "Cluster Pod Utilization",
                            "tooltip": {
                                "msResolution": false,
                                "shared": true,
                                "sort": 0,
                                "value_type": "individual"
                            },
                            "type": "graph",
                            "xaxis": {
                                "mode": "time",
                                "show": true,
                                "values": []
                            },
                            "yaxes": [
                                {
                                    "format": "short",
                                    "logBase": 1,
                                    "show": true
                                },
                                {
                                    "format": "short",
                                    "logBase": 1,
                                    "show": true
                                }
                            ]
                        },
                        {
                            "colorBackground": false,
                            "colorValue": false,
                            "colors": [
                                "rgba(50, 172, 45, 0.97)",
                                "rgba(237, 129, 40, 0.89)",
                                "rgba(245, 54, 54, 0.9)"
                            ],
                            "datasource": "${DS_PROMETHEUS}",
                            "editable": true,
                            "format": "percent",
                            "gauge": {
                                "maxValue": 100,
                                "minValue": 0,
                                "show": true,
                                "thresholdLabels": false,
                                "thresholdMarkers": true
                            },
                            "hideTimeOverride": false,
                            "id": 7,
                            "links": [],
                            "mappingType": 1,
                            "mappingTypes": [
                                {
                                    "name": "value to text",
                                    "value": 1
                                },
                                {
                                    "name": "range to text",
                                    "value": 2
                                }
                            ],
                            "maxDataPoints": 100,
                            "nullPointMode": "connected",
                            "postfix": "",
                            "postfixFontSize": "50%",
                            "prefix": "",
                            "prefixFontSize": "50%",
                            "rangeMaps": [
                                {
                                    "from": "null",
                                    "text": "N/A",
                                    "to": "null"
                                }
                            ],
                            "span": 3,
                            "sparkline": {
                                "fillColor": "rgba(31, 118, 189, 0.18)",
                                "full": false,
                                "lineColor": "rgb(31, 120, 193)",
                                "show": false
                            },
                            "targets": [
                                {
                                    "expr": "100 - (sum(kube_node_status_capacity_pods) - sum(kube_pod_info)) / sum(kube_node_status_capacity_pods) * 100",
                                    "format": "time_series",
                                    "intervalFactor": 2,
                                    "legendFormat": "",
                                    "refId": "A",
                                    "step": 60,
                                    "target": ""
                                }
                            ],
                            "thresholds": "80, 90",
                            "title": "Pod Utilization",
                            "transparent": false,
                            "type": "singlestat",
                            "valueFontSize": "80%",
                            "valueMaps": [
                                {
                                    "op": "=",
                                    "text": "N/A",
                                    "value": "null"
                                }
                            ],
                            "valueName": "current"
                        }
                    ],
                    "showTitle": false,
                    "title": "New Row",
                    "titleSize": "h6"
                }
            ],
            "schemaVersion": 14,
            "sharedCrosshair": false,
            "style": "dark",
            "tags": [],
            "templating": {
                "list": []
            },
            "time": {
                "from": "now-1h",
                "to": "now"
            },
            "timepicker": {
                "refresh_intervals": [
                    "5s",
                    "10s",
                    "30s",
                    "1m",
                    "5m",
                    "15m",
                    "30m",
                    "1h",
                    "2h",
                    "1d"
                ],
                "time_options": [
                    "5m",
                    "15m",
                    "1h",
                    "6h",
                    "12h",
                    "24h",
                    "2d",
                    "7d",
                    "30d"
                ]
            },
            "timezone": "browser",
            "title": "Kubernetes Capacity Planning",
            "version": 4
        },
        "inputs": [
            {
                "name": "DS_PROMETHEUS",
                "pluginId": "prometheus",
                "type": "datasource",
                "value": "prometheus"
            }
        ],
        "overwrite": true
    }
  kubernetes-cluster-monitoring-dashboard.json: |
    {
        "dashboard": {
            "__inputs": [
                {
                "name": "DS_PROMETHEUS",
                "label": "prometheus",
                "description": "",
                "type": "datasource",
                "pluginId": "prometheus",
                "pluginName": "Prometheus"
                }
            ],
            "__requires": [
                {
                "type": "grafana",
                "id": "grafana",
                "name": "Grafana",
                "version": "5.0.0"
                },
                {
                "type": "panel",
                "id": "graph",
                "name": "Graph",
                "version": "5.0.0"
                },
                {
                "type": "datasource",
                "id": "prometheus",
                "name": "Prometheus",
                "version": "5.0.0"
                },
                {
                "type": "panel",
                "id": "singlestat",
                "name": "Singlestat",
                "version": "5.0.0"
                }
            ],
            "annotations": {
                "list": [
                {
                    "builtIn": 1,
                    "datasource": "-- Grafana --",
                    "enable": true,
                    "hide": true,
                    "iconColor": "rgba(0, 211, 255, 1)",
                    "name": "Annotations & Alerts",
                    "type": "dashboard"
                }
                ]
            },
            "description": "Monitors Kubernetes cluster using Prometheus. Shows overall cluster CPU / Memory / Filesystem usage as well as individual pod, containers, systemd services statistics. Uses cAdvisor metrics only.",
            "editable": true,
            "gnetId": 1621,
            "graphTooltip": 0,
            "id": null,
            "iteration": 1555579102228,
            "links": [],
            "panels": [
                {
                "collapsed": false,
                "gridPos": {
                    "h": 1,
                    "w": 24,
                    "x": 0,
                    "y": 0
                },
                "id": 33,
                "panels": [],
                "title": "Network I/O pressure",
                "type": "row"
                },
                {
                "aliasColors": {},
                "bars": false,
                "dashLength": 10,
                "dashes": false,
                "datasource": "${DS_PROMETHEUS}",
                "decimals": 2,
                "editable": true,
                "error": false,
                "fill": 1,
                "grid": {},
                "gridPos": {
                    "h": 5,
                    "w": 24,
                    "x": 0,
                    "y": 1
                },
                "height": "200px",
                "id": 32,
                "isNew": true,
                "legend": {
                    "alignAsTable": false,
                    "avg": true,
                    "current": true,
                    "max": false,
                    "min": false,
                    "rightSide": false,
                    "show": false,
                    "sideWidth": 200,
                    "sort": "current",
                    "sortDesc": true,
                    "total": false,
                    "values": true
                },
                "lines": true,
                "linewidth": 2,
                "links": [],
                "nullPointMode": "connected",
                "percentage": false,
                "pointradius": 5,
                "points": false,
                "renderer": "flot",
                "seriesOverrides": [],
                "spaceLength": 10,
                "stack": false,
                "steppedLine": false,
                "targets": [
                    {
                    "expr": "sum (rate (container_network_receive_bytes_total{kubernetes_io_hostname=~\"^$Node$\"}[2m]))",
                    "interval": "10s",
                    "intervalFactor": 1,
                    "legendFormat": "Received",
                    "metric": "network",
                    "refId": "A",
                    "step": 10
                    },
                    {
                    "expr": "- sum (rate (container_network_transmit_bytes_total{kubernetes_io_hostname=~\"^$Node$\"}[2m]))",
                    "interval": "10s",
                    "intervalFactor": 1,
                    "legendFormat": "Sent",
                    "metric": "network",
                    "refId": "B",
                    "step": 10
                    }
                ],
                "thresholds": [],
                "timeFrom": null,
                "timeShift": null,
                "title": "Network I/O pressure",
                "tooltip": {
                    "msResolution": false,
                    "shared": true,
                    "sort": 0,
                    "value_type": "cumulative"
                },
                "transparent": false,
                "type": "graph",
                "xaxis": {
                    "buckets": null,
                    "mode": "time",
                    "name": null,
                    "show": true,
                    "values": []
                },
                "yaxes": [
                    {
                    "format": "Bps",
                    "label": null,
                    "logBase": 1,
                    "max": null,
                    "min": null,
                    "show": true
                    },
                    {
                    "format": "Bps",
                    "label": null,
                    "logBase": 1,
                    "max": null,
                    "min": null,
                    "show": false
                    }
                ]
                },
                {
                "collapsed": false,
                "gridPos": {
                    "h": 1,
                    "w": 24,
                    "x": 0,
                    "y": 6
                },
                "id": 34,
                "panels": [],
                "title": "Total usage",
                "type": "row"
                },
                {
                "cacheTimeout": null,
                "colorBackground": false,
                "colorValue": true,
                "colors": [
                    "rgba(50, 172, 45, 0.97)",
                    "rgba(237, 129, 40, 0.89)",
                    "rgba(245, 54, 54, 0.9)"
                ],
                "datasource": "${DS_PROMETHEUS}",
                "editable": true,
                "error": false,
                "format": "percent",
                "gauge": {
                    "maxValue": 100,
                    "minValue": 0,
                    "show": true,
                    "thresholdLabels": false,
                    "thresholdMarkers": true
                },
                "gridPos": {
                    "h": 5,
                    "w": 8,
                    "x": 0,
                    "y": 7
                },
                "height": "180px",
                "id": 4,
                "interval": null,
                "isNew": true,
                "links": [],
                "mappingType": 1,
                "mappingTypes": [
                    {
                    "name": "value to text",
                    "value": 1
                    },
                    {
                    "name": "range to text",
                    "value": 2
                    }
                ],
                "maxDataPoints": 100,
                "nullPointMode": "connected",
                "nullText": null,
                "postfix": "",
                "postfixFontSize": "50%",
                "prefix": "",
                "prefixFontSize": "50%",
                "rangeMaps": [
                    {
                    "from": "null",
                    "text": "N/A",
                    "to": "null"
                    }
                ],
                "sparkline": {
                    "fillColor": "rgba(31, 118, 189, 0.18)",
                    "full": false,
                    "lineColor": "rgb(31, 120, 193)",
                    "show": false
                },
                "tableColumn": "",
                "targets": [
                    {
                    "expr": "sum (container_memory_working_set_bytes{id=\"/\",kubernetes_io_hostname=~\"^$Node$\"}) / sum (machine_memory_bytes{kubernetes_io_hostname=~\"^$Node$\"}) * 100",
                    "interval": "10s",
                    "intervalFactor": 1,
                    "refId": "A",
                    "step": 10
                    }
                ],
                "thresholds": "65, 90",
                "title": "Cluster memory usage",
                "transparent": false,
                "type": "singlestat",
                "valueFontSize": "80%",
                "valueMaps": [
                    {
                    "op": "=",
                    "text": "N/A",
                    "value": "null"
                    }
                ],
                "valueName": "current"
                },
                {
                "cacheTimeout": null,
                "colorBackground": false,
                "colorValue": true,
                "colors": [
                    "rgba(50, 172, 45, 0.97)",
                    "rgba(237, 129, 40, 0.89)",
                    "rgba(245, 54, 54, 0.9)"
                ],
                "datasource": "${DS_PROMETHEUS}",
                "decimals": 2,
                "editable": true,
                "error": false,
                "format": "percent",
                "gauge": {
                    "maxValue": 100,
                    "minValue": 0,
                    "show": true,
                    "thresholdLabels": false,
                    "thresholdMarkers": true
                },
                "gridPos": {
                    "h": 5,
                    "w": 8,
                    "x": 8,
                    "y": 7
                },
                "height": "180px",
                "id": 6,
                "interval": null,
                "isNew": true,
                "links": [],
                "mappingType": 1,
                "mappingTypes": [
                    {
                    "name": "value to text",
                    "value": 1
                    },
                    {
                    "name": "range to text",
                    "value": 2
                    }
                ],
                "maxDataPoints": 100,
                "nullPointMode": "connected",
                "nullText": null,
                "postfix": "",
                "postfixFontSize": "50%",
                "prefix": "",
                "prefixFontSize": "50%",
                "rangeMaps": [
                    {
                    "from": "null",
                    "text": "N/A",
                    "to": "null"
                    }
                ],
                "sparkline": {
                    "fillColor": "rgba(31, 118, 189, 0.18)",
                    "full": false,
                    "lineColor": "rgb(31, 120, 193)",
                    "show": false
                },
                "tableColumn": "",
                "targets": [
                    {
                    "expr": "sum (rate (container_cpu_usage_seconds_total{id=\"/\",kubernetes_io_hostname=~\"^$Node$\"}[2m])) / sum (machine_cpu_cores{kubernetes_io_hostname=~\"^$Node$\"}) * 100",
                    "interval": "10s",
                    "intervalFactor": 1,
                    "refId": "A",
                    "step": 10
                    }
                ],
                "thresholds": "65, 90",
                "title": "Cluster CPU usage (2m avg)",
                "type": "singlestat",
                "valueFontSize": "80%",
                "valueMaps": [
                    {
                    "op": "=",
                    "text": "N/A",
                    "value": "null"
                    }
                ],
                "valueName": "current"
                },
                {
                "cacheTimeout": null,
                "colorBackground": false,
                "colorValue": true,
                "colors": [
                    "rgba(50, 172, 45, 0.97)",
                    "rgba(237, 129, 40, 0.89)",
                    "rgba(245, 54, 54, 0.9)"
                ],
                "datasource": "${DS_PROMETHEUS}",
                "decimals": 2,
                "editable": true,
                "error": false,
                "format": "percent",
                "gauge": {
                    "maxValue": 100,
                    "minValue": 0,
                    "show": true,
                    "thresholdLabels": false,
                    "thresholdMarkers": true
                },
                "gridPos": {
                    "h": 5,
                    "w": 8,
                    "x": 16,
                    "y": 7
                },
                "height": "180px",
                "id": 7,
                "interval": null,
                "isNew": true,
                "links": [],
                "mappingType": 1,
                "mappingTypes": [
                    {
                    "name": "value to text",
                    "value": 1
                    },
                    {
                    "name": "range to text",
                    "value": 2
                    }
                ],
                "maxDataPoints": 100,
                "nullPointMode": "connected",
                "nullText": null,
                "postfix": "",
                "postfixFontSize": "50%",
                "prefix": "",
                "prefixFontSize": "50%",
                "rangeMaps": [
                    {
                    "from": "null",
                    "text": "N/A",
                    "to": "null"
                    }
                ],
                "sparkline": {
                    "fillColor": "rgba(31, 118, 189, 0.18)",
                    "full": false,
                    "lineColor": "rgb(31, 120, 193)",
                    "show": false
                },
                "tableColumn": "",
                "targets": [
                    {
                    "expr": "sum (container_fs_usage_bytes{device=~\"^/dev/.*$\",id=\"/\",kubernetes_io_hostname=~\"^$Node$\"}) / sum (container_fs_limit_bytes{device=~\"^/dev/.*$\",id=\"/\",kubernetes_io_hostname=~\"^$Node$\"}) * 100",
                    "interval": "10s",
                    "intervalFactor": 1,
                    "legendFormat": "",
                    "metric": "",
                    "refId": "A",
                    "step": 10
                    }
                ],
                "thresholds": "65, 90",
                "title": "Cluster filesystem usage",
                "type": "singlestat",
                "valueFontSize": "80%",
                "valueMaps": [
                    {
                    "op": "=",
                    "text": "N/A",
                    "value": "null"
                    }
                ],
                "valueName": "current"
                },
                {
                "cacheTimeout": null,
                "colorBackground": false,
                "colorValue": false,
                "colors": [
                    "rgba(50, 172, 45, 0.97)",
                    "rgba(237, 129, 40, 0.89)",
                    "rgba(245, 54, 54, 0.9)"
                ],
                "datasource": "${DS_PROMETHEUS}",
                "decimals": 2,
                "editable": true,
                "error": false,
                "format": "bytes",
                "gauge": {
                    "maxValue": 100,
                    "minValue": 0,
                    "show": false,
                    "thresholdLabels": false,
                    "thresholdMarkers": true
                },
                "gridPos": {
                    "h": 3,
                    "w": 4,
                    "x": 0,
                    "y": 12
                },
                "height": "1px",
                "id": 9,
                "interval": null,
                "isNew": true,
                "links": [],
                "mappingType": 1,
                "mappingTypes": [
                    {
                    "name": "value to text",
                    "value": 1
                    },
                    {
                    "name": "range to text",
                    "value": 2
                    }
                ],
                "maxDataPoints": 100,
                "nullPointMode": "connected",
                "nullText": null,
                "postfix": "",
                "postfixFontSize": "20%",
                "prefix": "",
                "prefixFontSize": "20%",
                "rangeMaps": [
                    {
                    "from": "null",
                    "text": "N/A",
                    "to": "null"
                    }
                ],
                "sparkline": {
                    "fillColor": "rgba(31, 118, 189, 0.18)",
                    "full": false,
                    "lineColor": "rgb(31, 120, 193)",
                    "show": false
                },
                "tableColumn": "",
                "targets": [
                    {
                    "expr": "sum (container_memory_working_set_bytes{id=\"/\",kubernetes_io_hostname=~\"^$Node$\"})",
                    "interval": "10s",
                    "intervalFactor": 1,
                    "refId": "A",
                    "step": 10
                    }
                ],
                "thresholds": "",
                "title": "Used",
                "type": "singlestat",
                "valueFontSize": "50%",
                "valueMaps": [
                    {
                    "op": "=",
                    "text": "N/A",
                    "value": "null"
                    }
                ],
                "valueName": "current"
                },
                {
                "cacheTimeout": null,
                "colorBackground": false,
                "colorValue": false,
                "colors": [
                    "rgba(50, 172, 45, 0.97)",
                    "rgba(237, 129, 40, 0.89)",
                    "rgba(245, 54, 54, 0.9)"
                ],
                "datasource": "${DS_PROMETHEUS}",
                "decimals": 2,
                "editable": true,
                "error": false,
                "format": "bytes",
                "gauge": {
                    "maxValue": 100,
                    "minValue": 0,
                    "show": false,
                    "thresholdLabels": false,
                    "thresholdMarkers": true
                },
                "gridPos": {
                    "h": 3,
                    "w": 4,
                    "x": 4,
                    "y": 12
                },
                "height": "1px",
                "id": 10,
                "interval": null,
                "isNew": true,
                "links": [],
                "mappingType": 1,
                "mappingTypes": [
                    {
                    "name": "value to text",
                    "value": 1
                    },
                    {
                    "name": "range to text",
                    "value": 2
                    }
                ],
                "maxDataPoints": 100,
                "nullPointMode": "connected",
                "nullText": null,
                "postfix": "",
                "postfixFontSize": "50%",
                "prefix": "",
                "prefixFontSize": "50%",
                "rangeMaps": [
                    {
                    "from": "null",
                    "text": "N/A",
                    "to": "null"
                    }
                ],
                "sparkline": {
                    "fillColor": "rgba(31, 118, 189, 0.18)",
                    "full": false,
                    "lineColor": "rgb(31, 120, 193)",
                    "show": false
                },
                "tableColumn": "",
                "targets": [
                    {
                    "expr": "sum (machine_memory_bytes{kubernetes_io_hostname=~\"^$Node$\"})",
                    "interval": "10s",
                    "intervalFactor": 1,
                    "refId": "A",
                    "step": 10
                    }
                ],
                "thresholds": "",
                "title": "Total",
                "type": "singlestat",
                "valueFontSize": "50%",
                "valueMaps": [
                    {
                    "op": "=",
                    "text": "N/A",
                    "value": "null"
                    }
                ],
                "valueName": "current"
                },
                {
                "cacheTimeout": null,
                "colorBackground": false,
                "colorValue": false,
                "colors": [
                    "rgba(50, 172, 45, 0.97)",
                    "rgba(237, 129, 40, 0.89)",
                    "rgba(245, 54, 54, 0.9)"
                ],
                "datasource": "${DS_PROMETHEUS}",
                "decimals": 2,
                "editable": true,
                "error": false,
                "format": "none",
                "gauge": {
                    "maxValue": 100,
                    "minValue": 0,
                    "show": false,
                    "thresholdLabels": false,
                    "thresholdMarkers": true
                },
                "gridPos": {
                    "h": 3,
                    "w": 4,
                    "x": 8,
                    "y": 12
                },
                "height": "1px",
                "id": 11,
                "interval": null,
                "isNew": true,
                "links": [],
                "mappingType": 1,
                "mappingTypes": [
                    {
                    "name": "value to text",
                    "value": 1
                    },
                    {
                    "name": "range to text",
                    "value": 2
                    }
                ],
                "maxDataPoints": 100,
                "nullPointMode": "connected",
                "nullText": null,
                "postfix": " cores",
                "postfixFontSize": "30%",
                "prefix": "",
                "prefixFontSize": "50%",
                "rangeMaps": [
                    {
                    "from": "null",
                    "text": "N/A",
                    "to": "null"
                    }
                ],
                "sparkline": {
                    "fillColor": "rgba(31, 118, 189, 0.18)",
                    "full": false,
                    "lineColor": "rgb(31, 120, 193)",
                    "show": false
                },
                "tableColumn": "",
                "targets": [
                    {
                    "expr": "sum (rate (container_cpu_usage_seconds_total{id=\"/\",kubernetes_io_hostname=~\"^$Node$\"}[2m]))",
                    "interval": "10s",
                    "intervalFactor": 1,
                    "refId": "A",
                    "step": 10
                    }
                ],
                "thresholds": "",
                "title": "Used",
                "type": "singlestat",
                "valueFontSize": "50%",
                "valueMaps": [
                    {
                    "op": "=",
                    "text": "N/A",
                    "value": "null"
                    }
                ],
                "valueName": "current"
                },
                {
                "cacheTimeout": null,
                "colorBackground": false,
                "colorValue": false,
                "colors": [
                    "rgba(50, 172, 45, 0.97)",
                    "rgba(237, 129, 40, 0.89)",
                    "rgba(245, 54, 54, 0.9)"
                ],
                "datasource": "${DS_PROMETHEUS}",
                "decimals": 2,
                "editable": true,
                "error": false,
                "format": "none",
                "gauge": {
                    "maxValue": 100,
                    "minValue": 0,
                    "show": false,
                    "thresholdLabels": false,
                    "thresholdMarkers": true
                },
                "gridPos": {
                    "h": 3,
                    "w": 4,
                    "x": 12,
                    "y": 12
                },
                "height": "1px",
                "id": 12,
                "interval": null,
                "isNew": true,
                "links": [],
                "mappingType": 1,
                "mappingTypes": [
                    {
                    "name": "value to text",
                    "value": 1
                    },
                    {
                    "name": "range to text",
                    "value": 2
                    }
                ],
                "maxDataPoints": 100,
                "nullPointMode": "connected",
                "nullText": null,
                "postfix": " cores",
                "postfixFontSize": "30%",
                "prefix": "",
                "prefixFontSize": "50%",
                "rangeMaps": [
                    {
                    "from": "null",
                    "text": "N/A",
                    "to": "null"
                    }
                ],
                "sparkline": {
                    "fillColor": "rgba(31, 118, 189, 0.18)",
                    "full": false,
                    "lineColor": "rgb(31, 120, 193)",
                    "show": false
                },
                "tableColumn": "",
                "targets": [
                    {
                    "expr": "sum (machine_cpu_cores{kubernetes_io_hostname=~\"^$Node$\"})",
                    "interval": "10s",
                    "intervalFactor": 1,
                    "refId": "A",
                    "step": 10
                    }
                ],
                "thresholds": "",
                "title": "Total",
                "type": "singlestat",
                "valueFontSize": "50%",
                "valueMaps": [
                    {
                    "op": "=",
                    "text": "N/A",
                    "value": "null"
                    }
                ],
                "valueName": "current"
                },
                {
                "cacheTimeout": null,
                "colorBackground": false,
                "colorValue": false,
                "colors": [
                    "rgba(50, 172, 45, 0.97)",
                    "rgba(237, 129, 40, 0.89)",
                    "rgba(245, 54, 54, 0.9)"
                ],
                "datasource": "${DS_PROMETHEUS}",
                "decimals": 2,
                "editable": true,
                "error": false,
                "format": "bytes",
                "gauge": {
                    "maxValue": 100,
                    "minValue": 0,
                    "show": false,
                    "thresholdLabels": false,
                    "thresholdMarkers": true
                },
                "gridPos": {
                    "h": 3,
                    "w": 4,
                    "x": 16,
                    "y": 12
                },
                "height": "1px",
                "id": 13,
                "interval": null,
                "isNew": true,
                "links": [],
                "mappingType": 1,
                "mappingTypes": [
                    {
                    "name": "value to text",
                    "value": 1
                    },
                    {
                    "name": "range to text",
                    "value": 2
                    }
                ],
                "maxDataPoints": 100,
                "nullPointMode": "connected",
                "nullText": null,
                "postfix": "",
                "postfixFontSize": "50%",
                "prefix": "",
                "prefixFontSize": "50%",
                "rangeMaps": [
                    {
                    "from": "null",
                    "text": "N/A",
                    "to": "null"
                    }
                ],
                "sparkline": {
                    "fillColor": "rgba(31, 118, 189, 0.18)",
                    "full": false,
                    "lineColor": "rgb(31, 120, 193)",
                    "show": false
                },
                "tableColumn": "",
                "targets": [
                    {
                    "expr": "sum (container_fs_usage_bytes{device=~\"^/dev/.*$\",id=\"/\",kubernetes_io_hostname=~\"^$Node$\"})",
                    "interval": "10s",
                    "intervalFactor": 1,
                    "refId": "A",
                    "step": 10
                    }
                ],
                "thresholds": "",
                "title": "Used",
                "type": "singlestat",
                "valueFontSize": "50%",
                "valueMaps": [
                    {
                    "op": "=",
                    "text": "N/A",
                    "value": "null"
                    }
                ],
                "valueName": "current"
                },
                {
                "cacheTimeout": null,
                "colorBackground": false,
                "colorValue": false,
                "colors": [
                    "rgba(50, 172, 45, 0.97)",
                    "rgba(237, 129, 40, 0.89)",
                    "rgba(245, 54, 54, 0.9)"
                ],
                "datasource": "${DS_PROMETHEUS}",
                "decimals": 2,
                "editable": true,
                "error": false,
                "format": "bytes",
                "gauge": {
                    "maxValue": 100,
                    "minValue": 0,
                    "show": false,
                    "thresholdLabels": false,
                    "thresholdMarkers": true
                },
                "gridPos": {
                    "h": 3,
                    "w": 4,
                    "x": 20,
                    "y": 12
                },
                "height": "1px",
                "id": 14,
                "interval": null,
                "isNew": true,
                "links": [],
                "mappingType": 1,
                "mappingTypes": [
                    {
                    "name": "value to text",
                    "value": 1
                    },
                    {
                    "name": "range to text",
                    "value": 2
                    }
                ],
                "maxDataPoints": 100,
                "nullPointMode": "connected",
                "nullText": null,
                "postfix": "",
                "postfixFontSize": "50%",
                "prefix": "",
                "prefixFontSize": "50%",
                "rangeMaps": [
                    {
                    "from": "null",
                    "text": "N/A",
                    "to": "null"
                    }
                ],
                "sparkline": {
                    "fillColor": "rgba(31, 118, 189, 0.18)",
                    "full": false,
                    "lineColor": "rgb(31, 120, 193)",
                    "show": false
                },
                "tableColumn": "",
                "targets": [
                    {
                    "expr": "sum (container_fs_limit_bytes{device=~\"^/dev/.*$\",id=\"/\",kubernetes_io_hostname=~\"^$Node$\"})",
                    "interval": "10s",
                    "intervalFactor": 1,
                    "refId": "A",
                    "step": 10
                    }
                ],
                "thresholds": "",
                "title": "Total",
                "type": "singlestat",
                "valueFontSize": "50%",
                "valueMaps": [
                    {
                    "op": "=",
                    "text": "N/A",
                    "value": "null"
                    }
                ],
                "valueName": "current"
                }
            ],
            "refresh": "10s",
            "schemaVersion": 16,
            "style": "dark",
            "tags": [],
            "templating": {
                "list": [
                {
                    "allValue": ".*",
                    "current": {},
                    "datasource": "${DS_PROMETHEUS}",
                    "hide": 0,
                    "includeAll": true,
                    "label": null,
                    "multi": false,
                    "name": "Node",
                    "options": [],
                    "query": "label_values(kubernetes_io_hostname)",
                    "refresh": 1,
                    "regex": "",
                    "sort": 0,
                    "tagValuesQuery": "",
                    "tags": [],
                    "tagsQuery": "",
                    "type": "query",
                    "useTags": false
                }
                ]
            },
            "time": {
                "from": "now-5m",
                "to": "now"
            },
            "timepicker": {
                "refresh_intervals": [
                "5s",
                "10s",
                "30s",
                "1m",
                "5m",
                "15m",
                "30m",
                "1h",
                "2h",
                "1d"
                ],
                "time_options": [
                "5m",
                "15m",
                "1h",
                "6h",
                "12h",
                "24h",
                "2d",
                "7d",
                "30d"
                ]
            },
            "timezone": "browser",
            "title": "Kubernetes cluster monitoring",
            "uid": "EdZQa_gWz",
            "version": 3
        },
        "inputs": [
            {
                "name": "DS_PROMETHEUS",
                "pluginId": "prometheus",
                "type": "datasource",
                "value": "prometheus"
            }
        ],
        "overwrite": true
    }
  prometheus-datasource.json: |
    {
      "access": "proxy",
      "basicAuth": false,
      "name": "prometheus",
      "type": "prometheus",
      "url": "http://prometheus-server:80"
    }
`
