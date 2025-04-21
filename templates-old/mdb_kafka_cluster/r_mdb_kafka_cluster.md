---
subcategory: "Managed Service for Apache Kafka"
page_title: "Yandex: {{.Name}}"
description: |-
  Manages a Kafka cluster within Yandex Cloud.
---

# {{.Name}} ({{.Type}})

Manages a Kafka cluster within the Yandex Cloud. For more information, see [the official documentation](https://yandex.cloud/docs/managed-kafka/concepts).

## Example usage

{{ tffile "examples/mdb_kafka_cluster/r_mdb_kafka_cluster_1.tf" }}

Example of creating a HA Kafka Cluster with two brokers per AZ (6 brokers + 3 Zookepeers)

{{ tffile "examples/mdb_kafka_cluster/r_mdb_kafka_cluster_2.tf" }}

Example of creating Kafka Cluster with KRaft-controller subcluster instead of Zookeeper subcluster.

{{ tffile "examples/mdb_kafka_cluster/r_mdb_kafka_cluster_3.tf" }}

Example of creating multihost Kafka Cluster without subcluster of controllers, using KRaft-combine quorum.

{{ tffile "examples/mdb_kafka_cluster/r_mdb_kafka_cluster_4.tf" }}

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the Kafka cluster. Provided by the client when the cluster is created.

* `description` - (Optional) Description of the Kafka cluster.

* `folder_id` - (Optional) The ID of the folder that the resource belongs to. If it is not provided, the default provider folder is used.

* `labels` - (Optional) A set of key/value label pairs to assign to the Kafka cluster.

* `network_id` - (Required) ID of the network, to which the Kafka cluster belongs.

* `subnet_ids` - (Optional) IDs of the subnets, to which the Kafka cluster belongs.

* `environment` - (Optional) Deployment environment of the Kafka cluster. Can be either `PRESTABLE` or `PRODUCTION`. The default is `PRODUCTION`.

* `config` - (Required) Configuration of the Kafka cluster. The structure is documented below.

* `user` - (Deprecated) To manage users, please switch to using a separate resource type `yandex_mdb_kafka_user`.

* `topic` - (Deprecated) To manage topics, please switch to using a separate resource type `yandex_mdb_kafka_topic`.

* `security_group_ids` - (Optional) Security group ids, to which the Kafka cluster belongs.

* `host_group_ids` - (Optional) A list of IDs of the host groups to place VMs of the cluster on.

* `deletion_protection` - (Optional) Inhibits deletion of the cluster. Can be either `true` or `false`.

* `maintenance_window` - (Optional) Maintenance policy of the Kafka cluster. The structure is documented below.

~> Historically, `topic` blocks of the `yandex_mdb_kafka_cluster` resource were used to manage topics of the Kafka cluster. However, this approach has a number of disadvantages. In particular, when adding and removing topics from the tf recipe, terraform generates a diff that misleads the user about the planned changes. Also, this approach turned out to be inconvenient when managing topics through the Kafka Admin API. Therefore, topic management through a separate resource type `yandex_mdb_kafka_topic` was implemented and is now recommended.

---

The `maintenance_window` block supports:

* `type` - (Required) Type of maintenance window. Can be either `ANYTIME` or `WEEKLY`. A day and hour of window need to be specified with weekly window.

* `day` - (Optional) Day of the week (in `DDD` format). Allowed values: "MON", "TUE", "WED", "THU", "FRI", "SAT", "SUN"

* `hour` - (Optional) Hour of the day in UTC (in `HH` format). Allowed value is between 1 and 24.

---

The `config` block supports:

* `version` - (Required) Version of the Kafka server software.

* `brokers_count` - (Optional) Count of brokers per availability zone. The default is `1`.

* `zones` - (Required) List of availability zones.

* `assign_public_ip` - (Optional) Determines whether each broker will be assigned a public IP address. The default is `false`.

* `schema_registry` - (Optional) Enables managed schema registry on cluster. The default is `false`.

* `kafka` - (Optional) Configuration of the Kafka subcluster. The structure is documented below.

* `zookeeper` - (Optional) Configuration of the ZooKeeper subcluster. The structure is documented below.

* `kraft` - (Optional) Configuration of the KRaft-controller subcluster. The structure is documented below.

* `access` - (Optional) Access policy to the Kafka cluster. The structure is documented below.

* `disk_size_autoscaling` - (Optional) Disk autoscaling settings of the Kafka cluster. The structure is documented below.

---

The `kafka` block supports:

* `resources` - (Required) Resources allocated to hosts of the Kafka subcluster. The structure is documented below.

* `kafka_config` - (Optional) User-defined settings for the Kafka cluster. The structure is documented below.

The `resources` block supports:

* `resources_preset_id` - (Required) The ID of the preset for computational resources available to a Kafka host (CPU, memory etc.). For more information, see [the official documentation](https://yandex.cloud/docs/managed-kafka/concepts).

* `disk_size` - (Required) Volume of the storage available to a Kafka host, in gigabytes.

* `disk_type_id` - (Required) Type of the storage of Kafka hosts. For more information see [the official documentation](https://yandex.cloud/docs/managed-kafka/concepts/storage).

The `kafka_config` block supports:

* `compression_type`, `log_flush_interval_messages`, `log_flush_interval_ms`, `log_flush_scheduler_interval_ms`, `log_retention_bytes`, `log_retention_hours`, `log_retention_minutes`, `log_retention_ms`, `log_segment_bytes`, `socket_send_buffer_bytes`, `socket_receive_buffer_bytes`, `auto_create_topics_enable`, `num_partitions`, `default_replication_factor`, `message_max_bytes`, `replica_fetch_max_bytes`, `ssl_cipher_suites`, `offsets_retention_minutes`, `sasl_enabled_mechanisms` - (Optional) Kafka server settings. For more information, see [the official documentation](https://yandex.cloud/docs/managed-kafka/operations/cluster-update) and [the Kafka documentation](https://kafka.apache.org/documentation/#configuration).

The `zookeeper` block supports:

* `resources` - (Optional) Resources allocated to hosts of the ZooKeeper subcluster. The structure is documented below.

The `resources` block supports:

* `resources_preset_id` - (Optional) The ID of the preset for computational resources available to a ZooKeeper host (CPU, memory etc.). For more information, see [the official documentation](https://yandex.cloud/docs/managed-kafka/concepts).

* `disk_size` - (Optional) Volume of the storage available to a ZooKeeper host, in gigabytes.

* `disk_type_id` - (Optional) Type of the storage of ZooKeeper hosts. For more information see [the official documentation](https://yandex.cloud/docs/managed-kafka/concepts/storage).

The `kraft` block supports:

* `resources` - (Optional) Resources allocated to hosts of the KRaft-controller subcluster. The structure is documented below.

The `resources` block supports:

* `resources_preset_id` - (Optional) The ID of the preset for computational resources available to a KRaft-controller host (CPU, memory etc.). For more information, see [the official documentation](https://yandex.cloud/docs/managed-kafka/concepts).
* `disk_size` - (Optional) Volume of the storage available to a KRaft-controller host, in gigabytes.
* `disk_type_id` - (Optional) Type of the storage of KRaft-controller hosts. For more information see [the official documentation](https://yandex.cloud/docs/managed-kafka/concepts/storage).

The `access` block supports:

* `data_transfer` - Allow access for [DataTransfer](https://yandex.cloud/services/data-transfer)

The `disk_size_autoscaling` block supports:

* `disk_size_limit` - (Required) Maximum possible size of disk in bytes.

* `planned_usage_threshold` - (Optional) Percent of disk utilization. During maintenance disk will autoscale, if this threshold reached. Value is between 0 and 100. Default value is 0 (autoscaling disabled).

* `emergency_usage_threshold` - (Optional) Percent of disk utilization. Disk will autoscale immediately, if this threshold reached. Value is between 0 and 100. Default value is 0 (autoscaling disabled). Must be not less then 'planned_usage_threshold' value.

The `user` block is deprecated. To manage users, please switch to using a separate resource type `yandex_mdb_kafka_user`. The `user` block supports:

* `name` - (Required) The name of the user.

* `password` - (Required) The password of the user.

* `permission` - (Optional) Set of permissions granted to the user. The structure is documented below.

The `permission` block supports:

* `topic_name` - (Required) The name of the topic that the permission grants access to.

* `role` - (Required) The role type to grant to the topic.

* `allow_hosts` - (Optional) Set of hosts, to which this permission grants access to.

The `topic` block is deprecated. To manage topics, please switch to using a separate resource type `yandex_mdb_kafka_topic`. The `topic` block supports:

* `name` - (Required) The name of the topic.

* `partitions` - (Required) The number of the topic's partitions.

* `replication_factor` - (Required) Amount of data copies (replicas) for the topic in the cluster.

* `topic_config` - (Required) User-defined settings for the topic. The structure is documented below.

The `topic_config` block supports:

* `compression_type`, `delete_retention_ms`, `file_delete_delay_ms`, `flush_messages`, `flush_ms`, `min_compaction_lag_ms`, `retention_bytes`, `retention_ms`, `max_message_bytes`, `min_insync_replicas`, `segment_bytes`, `preallocate`, - (Optional) Kafka topic settings. For more information, see [the official documentation](https://yandex.cloud/docs/managed-kafka/operations/cluster-topics#update-topic) and [the Kafka documentation](https://kafka.apache.org/documentation/#configuration).

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `created_at` - Timestamp of cluster creation.

* `health` - Aggregated health of the cluster. Can be either `ALIVE`, `DEGRADED`, `DEAD` or `HEALTH_UNKNOWN`. For more information see `health` field of JSON representation in [the official documentation](https://yandex.cloud/docs/managed-kafka/api-ref/Cluster/).

* `status` - Status of the cluster. Can be either `CREATING`, `STARTING`, `RUNNING`, `UPDATING`, `STOPPING`, `STOPPED`, `ERROR` or `STATUS_UNKNOWN`. For more information see `status` field of JSON representation in [the official documentation](https://yandex.cloud/docs/managed-kafka/api-ref/Cluster/).

* `host` - A host of the Kafka cluster. The structure is documented below.

The `host` block supports:

* `name` - The fully qualified domain name of the host.
* `zone_id` - The availability zone where the Kafka host was created.
* `role` - Role of the host in the cluster.
* `health` - Health of the host.
* `subnet_id` - The ID of the subnet, to which the host belongs.
* `assign_public_ip` - The flag that defines whether a public IP address is assigned to the node.


## Import

The resource can be imported by using their `resource ID`. For getting the resource ID you can use Yandex Cloud [Web Console](https://console.yandex.cloud) or [YC CLI](https://yandex.cloud/docs/cli/quickstart).

{{ codefile "shell" "examples/mdb_kafka_cluster/import.sh" }}
