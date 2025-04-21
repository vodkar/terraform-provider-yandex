package yandex

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"

	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/mdb/kafka/v1"
	"github.com/yandex-cloud/terraform-provider-yandex/yandex/internal/hashcode"
)

type TopicCleanupPolicy int32

const (
	Topic_CLEANUP_POLICY_UNSPECIFIED TopicCleanupPolicy = 0
	// this policy discards log segments when either their retention time or log size limit is reached. See also: [KafkaConfig2_8.log_retention_ms] and other similar parameters.
	Topic_CLEANUP_POLICY_DELETE TopicCleanupPolicy = 1
	// this policy compacts messages in log.
	Topic_CLEANUP_POLICY_COMPACT TopicCleanupPolicy = 2
	// this policy use both compaction and deletion for messages and log segments.
	Topic_CLEANUP_POLICY_COMPACT_AND_DELETE TopicCleanupPolicy = 3
)

const kafkaConfigPath = "config.0.kafka.0.kafka_config.0"

// Enum value maps for TopicCleanupPolicy.
var (
	Topic_CleanupPolicy_name = map[int32]string{
		0: "CLEANUP_POLICY_UNSPECIFIED",
		1: "CLEANUP_POLICY_DELETE",
		2: "CLEANUP_POLICY_COMPACT",
		3: "CLEANUP_POLICY_COMPACT_AND_DELETE",
	}
	Topic_CleanupPolicy_value = map[string]int32{
		"CLEANUP_POLICY_UNSPECIFIED":        0,
		"CLEANUP_POLICY_DELETE":             1,
		"CLEANUP_POLICY_COMPACT":            2,
		"CLEANUP_POLICY_COMPACT_AND_DELETE": 3,
	}
)

func parseKafkaEnv(e string) (kafka.Cluster_Environment, error) {
	v, ok := kafka.Cluster_Environment_value[e]
	if !ok {
		return 0, fmt.Errorf("value for 'environment' must be one of %s, not `%s`",
			getJoinedKeys(getEnumValueMapKeys(kafka.Cluster_Environment_value)), e)
	}
	return kafka.Cluster_Environment(v), nil
}

func parseKafkaCompression(e string) (kafka.CompressionType, error) {
	v, ok := kafka.CompressionType_value[e]
	if !ok || e == "COMPRESSION_TYPE_UNSPECIFIED" {
		return 0, fmt.Errorf("value for 'compression_type' must be one of %s, not `%s`",
			getJoinedKeys(getEnumValueMapKeysExt(kafka.CompressionType_value, true)), e)
	}
	return kafka.CompressionType(v), nil
}

func parseKafkaSaslMechanism(e string) (kafka.SaslMechanism, error) {
	v, ok := kafka.SaslMechanism_value[e]
	if !ok || e == "SASL_MECHANISM_UNSPECIFIED" {
		return 0, fmt.Errorf("value for 'sasl_mechanism' must be one of %s, not `%s`",
			getJoinedKeys(getEnumValueMapKeysExt(kafka.SaslMechanism_value, true)), e)
	}
	return kafka.SaslMechanism(v), nil
}

func parseKafkaPermissionRole(e string) (kafka.Permission_AccessRole, error) {
	v, ok := kafka.Permission_AccessRole_value[e]
	if !ok {
		return 0, fmt.Errorf("value for 'role' must be one of %s, not `%s`",
			getJoinedKeys(getEnumValueMapKeys(kafka.Permission_AccessRole_value)), e)
	}
	return kafka.Permission_AccessRole(v), nil
}

func parseSetToStringArray(set interface{}) []string {
	if set == nil {
		return nil
	}
	schemaSet := set.(*schema.Set)
	return convertStringSet(schemaSet)
}

func parseKafkaPermissionAllowHosts(allowHosts interface{}) []string {
	result := parseSetToStringArray(allowHosts)
	if len(result) == 0 {
		return nil
	}
	sort.Strings(result)
	return result
}

func parseKafkaTopicCleanupPolicy(e string) (TopicCleanupPolicy, error) {
	v, ok := Topic_CleanupPolicy_value[e]
	if !ok || e == "CLEANUP_POLICY_UNSPECIFIED" {
		return 0, fmt.Errorf("value for 'cleanup_policy' must be one of %s, not `%s`",
			getJoinedKeys(getEnumValueMapKeysExt(Topic_CleanupPolicy_value, true)), e)
	}
	return TopicCleanupPolicy(v), nil
}

func parseIntKafkaConfigParam(d *schema.ResourceData, paramName string, retErr *error) *wrappers.Int64Value {
	v, ok := d.GetOk(kafkaConfigPath + "." + paramName)
	if !ok {
		return nil
	}

	i, err := strconv.ParseInt(v.(string), 10, 64)
	if err != nil {
		if *retErr != nil {
			*retErr = err
		}
		return nil
	}
	return &wrappers.Int64Value{Value: i}
}

func parseSslCipherSuites(sslCipherSuites interface{}) []string {
	return parseSetToStringArray(sslCipherSuites)
}

func parseSaslEnabledMechanisms(saslEnabledMechanisms interface{}) ([]kafka.SaslMechanism, error) {
	if saslEnabledMechanisms == nil {
		return nil, nil
	}
	setOfMechanisms := saslEnabledMechanisms.(*schema.Set)
	sliceOfMechanisms := convertStringSet(setOfMechanisms)
	var result []kafka.SaslMechanism
	for _, mechanismString := range sliceOfMechanisms {
		mechanismSpec, err := parseKafkaSaslMechanism(mechanismString)
		if err != nil {
			return nil, err
		}
		result = append(result, mechanismSpec)
	}
	return result, nil
}

type KafkaConfig struct {
	CompressionType             kafka.CompressionType
	LogFlushIntervalMessages    *wrappers.Int64Value
	LogFlushIntervalMs          *wrappers.Int64Value
	LogFlushSchedulerIntervalMs *wrappers.Int64Value
	LogRetentionBytes           *wrappers.Int64Value
	LogRetentionHours           *wrappers.Int64Value
	LogRetentionMinutes         *wrappers.Int64Value
	LogRetentionMs              *wrappers.Int64Value
	LogSegmentBytes             *wrappers.Int64Value
	LogPreallocate              *wrappers.BoolValue
	SocketSendBufferBytes       *wrappers.Int64Value
	SocketReceiveBufferBytes    *wrappers.Int64Value
	AutoCreateTopicsEnable      *wrappers.BoolValue
	NumPartitions               *wrappers.Int64Value
	DefaultReplicationFactor    *wrappers.Int64Value
	MessageMaxBytes             *wrappers.Int64Value
	ReplicaFetchMaxBytes        *wrappers.Int64Value
	SslCipherSuites             []string
	OffsetsRetentionMinutes     *wrappers.Int64Value
	SaslEnabledMechanisms       []kafka.SaslMechanism
}

func parseKafkaConfig(d *schema.ResourceData) (*KafkaConfig, error) {
	res := &KafkaConfig{}

	if v, ok := d.GetOk(kafkaConfigPath + ".compression_type"); ok {
		value, err := parseKafkaCompression(v.(string))
		if err != nil {
			return nil, err
		}
		res.CompressionType = value
	}

	var retErr error

	res.LogFlushIntervalMessages = parseIntKafkaConfigParam(d, "log_flush_interval_messages", &retErr)
	res.LogFlushIntervalMs = parseIntKafkaConfigParam(d, "log_flush_interval_ms", &retErr)
	res.LogFlushSchedulerIntervalMs = parseIntKafkaConfigParam(d, "log_flush_scheduler_interval_ms", &retErr)
	res.LogRetentionBytes = parseIntKafkaConfigParam(d, "log_retention_bytes", &retErr)
	res.LogRetentionHours = parseIntKafkaConfigParam(d, "log_retention_hours", &retErr)
	res.LogRetentionMinutes = parseIntKafkaConfigParam(d, "log_retention_minutes", &retErr)
	res.LogRetentionMs = parseIntKafkaConfigParam(d, "log_retention_ms", &retErr)
	res.LogSegmentBytes = parseIntKafkaConfigParam(d, "log_segment_bytes", &retErr)
	res.SocketSendBufferBytes = parseIntKafkaConfigParam(d, "socket_send_buffer_bytes", &retErr)
	res.SocketReceiveBufferBytes = parseIntKafkaConfigParam(d, "socket_receive_buffer_bytes", &retErr)
	res.NumPartitions = parseIntKafkaConfigParam(d, "num_partitions", &retErr)
	res.DefaultReplicationFactor = parseIntKafkaConfigParam(d, "default_replication_factor", &retErr)
	res.MessageMaxBytes = parseIntKafkaConfigParam(d, "message_max_bytes", &retErr)
	res.ReplicaFetchMaxBytes = parseIntKafkaConfigParam(d, "replica_fetch_max_bytes", &retErr)
	res.OffsetsRetentionMinutes = parseIntKafkaConfigParam(d, "offsets_retention_minutes", &retErr)

	// if v, ok := d.GetOk(kafkaConfigPath + ".log_preallocate"); ok {
	// 	res.LogPreallocate = &wrappers.BoolValue{Value: v.(bool)}
	// }
	if v, ok := d.GetOk(kafkaConfigPath + ".auto_create_topics_enable"); ok {
		res.AutoCreateTopicsEnable = &wrappers.BoolValue{Value: v.(bool)}
	}
	if v, ok := d.GetOk(kafkaConfigPath + ".ssl_cipher_suites"); ok {
		res.SslCipherSuites = parseSslCipherSuites(v)
		sort.Strings(res.SslCipherSuites)
	}
	if v, ok := d.GetOk(kafkaConfigPath + ".sasl_enabled_mechanisms"); ok {
		mechanisms, err := parseSaslEnabledMechanisms(v)
		if err != nil {
			return nil, err
		}
		res.SaslEnabledMechanisms = mechanisms
	}

	if retErr != nil {
		return nil, retErr
	}

	return res, nil
}

func expandKafkaConfig2_8(d *schema.ResourceData) (*kafka.KafkaConfig2_8, error) {
	kafkaConfig, err := parseKafkaConfig(d)
	if err != nil {
		return nil, err
	}
	return &kafka.KafkaConfig2_8{
		CompressionType:             kafkaConfig.CompressionType,
		LogFlushIntervalMessages:    kafkaConfig.LogFlushIntervalMessages,
		LogFlushIntervalMs:          kafkaConfig.LogFlushIntervalMs,
		LogFlushSchedulerIntervalMs: kafkaConfig.LogFlushSchedulerIntervalMs,
		LogRetentionBytes:           kafkaConfig.LogRetentionBytes,
		LogRetentionHours:           kafkaConfig.LogRetentionHours,
		LogRetentionMinutes:         kafkaConfig.LogRetentionMinutes,
		LogRetentionMs:              kafkaConfig.LogRetentionMs,
		LogSegmentBytes:             kafkaConfig.LogSegmentBytes,
		LogPreallocate:              kafkaConfig.LogPreallocate,
		SocketSendBufferBytes:       kafkaConfig.SocketSendBufferBytes,
		SocketReceiveBufferBytes:    kafkaConfig.SocketReceiveBufferBytes,
		AutoCreateTopicsEnable:      kafkaConfig.AutoCreateTopicsEnable,
		NumPartitions:               kafkaConfig.NumPartitions,
		DefaultReplicationFactor:    kafkaConfig.DefaultReplicationFactor,
		MessageMaxBytes:             kafkaConfig.MessageMaxBytes,
		ReplicaFetchMaxBytes:        kafkaConfig.ReplicaFetchMaxBytes,
		SslCipherSuites:             kafkaConfig.SslCipherSuites,
		OffsetsRetentionMinutes:     kafkaConfig.OffsetsRetentionMinutes,
		SaslEnabledMechanisms:       kafkaConfig.SaslEnabledMechanisms,
	}, nil
}

func expandKafkaConfig3x(d *schema.ResourceData) (*kafka.KafkaConfig3, error) {
	kafkaConfig, err := parseKafkaConfig(d)
	if err != nil {
		return nil, err
	}
	return &kafka.KafkaConfig3{
		CompressionType:             kafkaConfig.CompressionType,
		LogFlushIntervalMessages:    kafkaConfig.LogFlushIntervalMessages,
		LogFlushIntervalMs:          kafkaConfig.LogFlushIntervalMs,
		LogFlushSchedulerIntervalMs: kafkaConfig.LogFlushSchedulerIntervalMs,
		LogRetentionBytes:           kafkaConfig.LogRetentionBytes,
		LogRetentionHours:           kafkaConfig.LogRetentionHours,
		LogRetentionMinutes:         kafkaConfig.LogRetentionMinutes,
		LogRetentionMs:              kafkaConfig.LogRetentionMs,
		LogSegmentBytes:             kafkaConfig.LogSegmentBytes,
		LogPreallocate:              kafkaConfig.LogPreallocate,
		SocketSendBufferBytes:       kafkaConfig.SocketSendBufferBytes,
		SocketReceiveBufferBytes:    kafkaConfig.SocketReceiveBufferBytes,
		AutoCreateTopicsEnable:      kafkaConfig.AutoCreateTopicsEnable,
		NumPartitions:               kafkaConfig.NumPartitions,
		DefaultReplicationFactor:    kafkaConfig.DefaultReplicationFactor,
		MessageMaxBytes:             kafkaConfig.MessageMaxBytes,
		ReplicaFetchMaxBytes:        kafkaConfig.ReplicaFetchMaxBytes,
		SslCipherSuites:             kafkaConfig.SslCipherSuites,
		OffsetsRetentionMinutes:     kafkaConfig.OffsetsRetentionMinutes,
		SaslEnabledMechanisms:       kafkaConfig.SaslEnabledMechanisms,
	}, nil
}

type TopicConfig struct {
	CleanupPolicy      string
	CompressionType    kafka.CompressionType
	DeleteRetentionMs  *wrappers.Int64Value
	FileDeleteDelayMs  *wrappers.Int64Value
	FlushMessages      *wrappers.Int64Value
	FlushMs            *wrappers.Int64Value
	MinCompactionLagMs *wrappers.Int64Value
	RetentionBytes     *wrappers.Int64Value
	RetentionMs        *wrappers.Int64Value
	MaxMessageBytes    *wrappers.Int64Value
	MinInsyncReplicas  *wrappers.Int64Value
	SegmentBytes       *wrappers.Int64Value
	Preallocate        *wrappers.BoolValue
}

func parseIntTopicConfigParam(d *schema.ResourceData, paramPath string, retErr *error) *wrappers.Int64Value {
	paramValue, ok := d.GetOk(paramPath)
	if !ok {
		return nil
	}
	str := paramValue.(string)
	if str == "" {
		return nil
	}
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		if *retErr != nil {
			*retErr = err
		}
		return nil
	}
	return &wrappers.Int64Value{Value: i}
}

func parseKafkaTopicConfig(d *schema.ResourceData, topicConfigPrefix string) (*TopicConfig, error) {
	key := func(key string) string {
		return fmt.Sprintf("%s%s", topicConfigPrefix, key)
	}

	res := &TopicConfig{}
	if cleanupPolicy := d.Get(key("cleanup_policy")).(string); cleanupPolicy != "" {
		_, err := parseKafkaTopicCleanupPolicy(cleanupPolicy)
		if err != nil {
			return nil, err
		}
		res.CleanupPolicy = cleanupPolicy
	}

	if compressionType := d.Get(key("compression_type")).(string); compressionType != "" {
		value, err := parseKafkaCompression(compressionType)
		if err != nil {
			return nil, err
		}
		res.CompressionType = value
	}

	var retErr error
	res.DeleteRetentionMs = parseIntTopicConfigParam(d, key("delete_retention_ms"), &retErr)
	res.FileDeleteDelayMs = parseIntTopicConfigParam(d, key("file_delete_delay_ms"), &retErr)
	res.FlushMessages = parseIntTopicConfigParam(d, key("flush_messages"), &retErr)
	res.FlushMs = parseIntTopicConfigParam(d, key("flush_ms"), &retErr)
	res.MinCompactionLagMs = parseIntTopicConfigParam(d, key("min_compaction_lag_ms"), &retErr)
	res.RetentionBytes = parseIntTopicConfigParam(d, key("retention_bytes"), &retErr)
	res.RetentionMs = parseIntTopicConfigParam(d, key("retention_ms"), &retErr)
	res.MaxMessageBytes = parseIntTopicConfigParam(d, key("max_message_bytes"), &retErr)
	res.MinInsyncReplicas = parseIntTopicConfigParam(d, key("min_insync_replicas"), &retErr)
	res.SegmentBytes = parseIntTopicConfigParam(d, key("segment_bytes"), &retErr)

	if preallocateRaw, ok := d.GetOk(key("preallocate")); ok {
		res.Preallocate = &wrappers.BoolValue{Value: preallocateRaw.(bool)}
	}

	if retErr != nil {
		return nil, retErr
	}

	return res, nil
}

func expandKafkaTopicConfig2_8(d *schema.ResourceData, topicConfigPrefix string) (*kafka.TopicConfig2_8, error) {
	topicConfig, err := parseKafkaTopicConfig(d, topicConfigPrefix)
	if err != nil {
		return nil, err
	}
	res := &kafka.TopicConfig2_8{
		CleanupPolicy:      kafka.TopicConfig2_8_CleanupPolicy(kafka.TopicConfig2_8_CleanupPolicy_value[topicConfig.CleanupPolicy]),
		CompressionType:    topicConfig.CompressionType,
		DeleteRetentionMs:  topicConfig.DeleteRetentionMs,
		FileDeleteDelayMs:  topicConfig.FileDeleteDelayMs,
		FlushMessages:      topicConfig.FlushMessages,
		FlushMs:            topicConfig.FlushMs,
		MinCompactionLagMs: topicConfig.MinCompactionLagMs,
		RetentionBytes:     topicConfig.RetentionBytes,
		RetentionMs:        topicConfig.RetentionMs,
		MaxMessageBytes:    topicConfig.MaxMessageBytes,
		MinInsyncReplicas:  topicConfig.MinInsyncReplicas,
		SegmentBytes:       topicConfig.SegmentBytes,
		Preallocate:        topicConfig.Preallocate,
	}

	return res, nil
}

func expandKafkaTopicConfig3x(d *schema.ResourceData, topicConfigPrefix string) (*kafka.TopicConfig3, error) {
	topicConfig, err := parseKafkaTopicConfig(d, topicConfigPrefix)
	if err != nil {
		return nil, err
	}
	res := &kafka.TopicConfig3{
		CleanupPolicy:      kafka.TopicConfig3_CleanupPolicy(kafka.TopicConfig3_CleanupPolicy_value[topicConfig.CleanupPolicy]),
		CompressionType:    topicConfig.CompressionType,
		DeleteRetentionMs:  topicConfig.DeleteRetentionMs,
		FileDeleteDelayMs:  topicConfig.FileDeleteDelayMs,
		FlushMessages:      topicConfig.FlushMessages,
		FlushMs:            topicConfig.FlushMs,
		MinCompactionLagMs: topicConfig.MinCompactionLagMs,
		RetentionBytes:     topicConfig.RetentionBytes,
		RetentionMs:        topicConfig.RetentionMs,
		MaxMessageBytes:    topicConfig.MaxMessageBytes,
		MinInsyncReplicas:  topicConfig.MinInsyncReplicas,
		SegmentBytes:       topicConfig.SegmentBytes,
		Preallocate:        topicConfig.Preallocate,
	}

	return res, nil
}

func expandKafkaConfigSpec(d *schema.ResourceData) (*kafka.ConfigSpec, error) {
	result := &kafka.ConfigSpec{}

	if v, ok := d.GetOk("config.0.version"); ok {
		result.Version = v.(string)
	}

	if v, ok := d.GetOk("config.0.brokers_count"); ok {
		result.BrokersCount = &wrappers.Int64Value{Value: int64(v.(int))}
	}

	if v, ok := d.GetOk("config.0.assign_public_ip"); ok {
		result.AssignPublicIp = v.(bool)
	}

	if v, ok := d.GetOk("config.0.schema_registry"); ok {
		result.SchemaRegistry = v.(bool)
	}

	if v, ok := d.GetOk("config.0.zones"); ok {
		zones := v.([]interface{})
		result.ZoneId = []string{}
		for _, zone := range zones {
			result.ZoneId = append(result.ZoneId, zone.(string))
		}
	}
	result.Kafka = &kafka.ConfigSpec_Kafka{}
	result.Kafka.Resources = expandKafkaResources(d, "config.0.kafka.0.resources.0")

	version := result.Version
	if strings.HasPrefix(version, "3") {
		cfg, err := expandKafkaConfig3x(d)
		if err != nil {
			return nil, err
		}
		result.Kafka.SetKafkaConfig_3(cfg)
	} else if version == "2.8" {
		cfg, err := expandKafkaConfig2_8(d)
		if err != nil {
			return nil, err
		}
		result.Kafka.SetKafkaConfig_2_8(cfg)
	} else if version == "" {
		return nil, fmt.Errorf("you must specify version of Kafka")
	} else {
		return nil, fmt.Errorf("this version of Kafka not supported by Terraform provider")
	}

	if _, ok := d.GetOk("config.0.zookeeper"); ok {
		result.Zookeeper = &kafka.ConfigSpec_Zookeeper{}
		result.Zookeeper.Resources = expandKafkaResources(d, "config.0.zookeeper.0.resources.0")
	}

	if _, ok := d.GetOk("config.0.kraft"); ok {
		result.Kraft = &kafka.ConfigSpec_KRaft{}
		result.Kraft.Resources = expandKafkaResources(d, "config.0.kraft.0.resources.0")
	}

	result.SetAccess(expandKafkaAccess(d))
	result.SetRestApiConfig(expandKafkaRestAPI(d))
	result.DiskSizeAutoscaling = expandKafkaDiskSizeAutoscaling(d)

	return result, nil
}

func expandKafkaDiskSizeAutoscaling(d *schema.ResourceData) *kafka.DiskSizeAutoscaling {
	if _, ok := d.GetOkExists("config.0.disk_size_autoscaling"); !ok {
		return nil
	}

	out := &kafka.DiskSizeAutoscaling{}

	if v, ok := d.GetOk("config.0.disk_size_autoscaling.0.disk_size_limit"); ok {
		out.DiskSizeLimit = toBytes(v.(int))
	}

	if v, ok := d.GetOk("config.0.disk_size_autoscaling.0.planned_usage_threshold"); ok {
		out.PlannedUsageThreshold = int64(v.(int))
	}

	if v, ok := d.GetOk("config.0.disk_size_autoscaling.0.emergency_usage_threshold"); ok {
		out.EmergencyUsageThreshold = int64(v.(int))
	}

	return out
}

func expandKafkaTopics(d *schema.ResourceData) ([]*kafka.TopicSpec, error) {
	var result []*kafka.TopicSpec
	version, ok := d.GetOk("config.0.version")
	if !ok {
		return nil, fmt.Errorf("you must specify version of Kafka")
	}
	topics := d.Get("topic").([]interface{})

	for idx := range topics {
		topicSpec, err := buildKafkaTopicSpec(d, fmt.Sprintf("topic.%d.", idx), version.(string))
		if err != nil {
			return nil, err
		}
		result = append(result, topicSpec)
	}
	return result, nil
}

func expandKafkaUsers(d *schema.ResourceData) ([]*kafka.UserSpec, error) {
	users := d.Get("user").(*schema.Set)
	result := make([]*kafka.UserSpec, 0, users.Len())

	for _, u := range users.List() {
		user, err := expandKafkaUser(u)
		if err != nil {
			return nil, err
		}
		result = append(result, user)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result, nil
}

func expandKafkaUser(u interface{}) (*kafka.UserSpec, error) {
	m := u.(map[string]interface{})
	user := &kafka.UserSpec{}
	if v, ok := m["name"]; ok {
		user.Name = v.(string)
	}
	if v, ok := m["password"]; ok {
		user.Password = v.(string)
	}
	if v, ok := m["permission"]; ok {
		permissions, err := expandKafkaPermissions(v.(*schema.Set))
		if err != nil {
			return nil, err
		}
		user.Permissions = permissions
	}
	return user, nil
}

func expandKafkaPermissions(ps *schema.Set) ([]*kafka.Permission, error) {
	result := []*kafka.Permission{}

	for _, p := range ps.List() {
		m := p.(map[string]interface{})
		permission := &kafka.Permission{}
		if v, ok := m["topic_name"]; ok {
			permission.TopicName = v.(string)
		}
		if v, ok := m["role"]; ok {
			role, err := parseKafkaPermissionRole(v.(string))
			if err != nil {
				return nil, err
			}
			permission.Role = role
		}
		if v, ok := m["allow_hosts"]; ok {
			permission.AllowHosts = parseKafkaPermissionAllowHosts(v)
		}
		result = append(result, permission)
	}
	sortPermissions(result)
	return result, nil
}

func flattenKafkaConfig(cluster *kafka.Cluster) ([]map[string]interface{}, error) {
	kafkaResources, err := flattenKafkaResources(cluster.Config.Kafka.Resources)
	if err != nil {
		return nil, err
	}

	var kafkaConfig map[string]interface{}
	if cluster.Config.Kafka.GetKafkaConfig_2_8() != nil {
		kafkaConfig, err = flattenKafkaConfig2_8Settings(cluster.Config.Kafka.GetKafkaConfig_2_8())
		if err != nil {
			return nil, err
		}
	}
	if cluster.Config.Kafka.GetKafkaConfig_3() != nil {
		kafkaConfig, err = flattenKafkaConfig3Settings(cluster.Config.Kafka.GetKafkaConfig_3())
		if err != nil {
			return nil, err
		}
	}

	config := map[string]interface{}{
		"brokers_count":    cluster.Config.BrokersCount.GetValue(),
		"assign_public_ip": cluster.Config.AssignPublicIp,
		"schema_registry":  cluster.Config.SchemaRegistry,
		"zones":            cluster.Config.ZoneId,
		"version":          cluster.Config.Version,
		"kafka": []map[string]interface{}{
			{
				"resources":    []map[string]interface{}{kafkaResources},
				"kafka_config": []map[string]interface{}{kafkaConfig},
			},
		},
	}
	if cluster.Config.Zookeeper != nil {
		zkResources, err := flattenKafkaResources(cluster.Config.Zookeeper.Resources)
		if err != nil {
			return nil, err
		}
		config["zookeeper"] = []map[string]interface{}{
			{
				"resources": []map[string]interface{}{zkResources},
			},
		}
	}
	if cluster.Config.Kraft != nil {
		kRaftResources, err := flattenKafkaResources(cluster.Config.Kraft.Resources)
		if err != nil {
			return nil, err
		}
		config["kraft"] = []map[string]interface{}{
			{
				"resources": []map[string]interface{}{kRaftResources},
			},
		}
	}
	if cluster.Config.GetAccess() != nil {
		config["access"] = flattenKafkaAccess(cluster.Config)
	}
	if cluster.Config.GetRestApiConfig() != nil {
		config["rest_api"] = flattenKafkaRestAPI(cluster.Config)
	}
	config["disk_size_autoscaling"] = flattenKafkaDiskSizeAutoscaling(cluster.Config.DiskSizeAutoscaling)

	return []map[string]interface{}{config}, nil
}

func flattenKafkaDiskSizeAutoscaling(p *kafka.DiskSizeAutoscaling) []interface{} {
	if p == nil {
		return nil
	}

	out := map[string]interface{}{}

	out["disk_size_limit"] = toGigabytes(p.DiskSizeLimit)
	out["planned_usage_threshold"] = int(p.PlannedUsageThreshold)
	out["emergency_usage_threshold"] = int(p.EmergencyUsageThreshold)

	return []interface{}{out}
}

type KafkaConfigSettings interface {
	GetCompressionType() kafka.CompressionType
	GetLogFlushIntervalMessages() *wrappers.Int64Value
	GetLogFlushIntervalMs() *wrappers.Int64Value
	GetLogFlushSchedulerIntervalMs() *wrappers.Int64Value
	GetLogRetentionBytes() *wrappers.Int64Value
	GetLogRetentionHours() *wrappers.Int64Value
	GetLogRetentionMinutes() *wrappers.Int64Value
	GetLogRetentionMs() *wrappers.Int64Value
	GetLogSegmentBytes() *wrappers.Int64Value
	GetLogPreallocate() *wrappers.BoolValue
	GetSocketSendBufferBytes() *wrappers.Int64Value
	GetSocketReceiveBufferBytes() *wrappers.Int64Value
	GetAutoCreateTopicsEnable() *wrappers.BoolValue
	GetNumPartitions() *wrappers.Int64Value
	GetDefaultReplicationFactor() *wrappers.Int64Value
	GetMessageMaxBytes() *wrappers.Int64Value
	GetReplicaFetchMaxBytes() *wrappers.Int64Value
	GetSslCipherSuites() []string
	GetOffsetsRetentionMinutes() *wrappers.Int64Value
	GetSaslEnabledMechanisms() []kafka.SaslMechanism
}

func flattenKafkaConfigSettings(kafkaConfig KafkaConfigSettings) (map[string]interface{}, error) {
	res := map[string]interface{}{}

	if kafkaConfig.GetCompressionType() != kafka.CompressionType_COMPRESSION_TYPE_UNSPECIFIED {
		res["compression_type"] = kafkaConfig.GetCompressionType().String()
	}
	if kafkaConfig.GetLogFlushIntervalMessages() != nil {
		res["log_flush_interval_messages"] = strconv.FormatInt(kafkaConfig.GetLogFlushIntervalMessages().GetValue(), 10)
	}
	if kafkaConfig.GetLogFlushIntervalMs() != nil {
		res["log_flush_interval_ms"] = strconv.FormatInt(kafkaConfig.GetLogFlushIntervalMs().GetValue(), 10)
	}
	if kafkaConfig.GetLogFlushSchedulerIntervalMs() != nil {
		res["log_flush_scheduler_interval_ms"] = strconv.FormatInt(kafkaConfig.GetLogFlushSchedulerIntervalMs().GetValue(), 10)
	}
	if kafkaConfig.GetLogRetentionBytes() != nil {
		res["log_retention_bytes"] = strconv.FormatInt(kafkaConfig.GetLogRetentionBytes().GetValue(), 10)
	}
	if kafkaConfig.GetLogRetentionHours() != nil {
		res["log_retention_hours"] = strconv.FormatInt(kafkaConfig.GetLogRetentionHours().GetValue(), 10)
	}
	if kafkaConfig.GetLogRetentionMinutes() != nil {
		res["log_retention_minutes"] = strconv.FormatInt(kafkaConfig.GetLogRetentionMinutes().GetValue(), 10)
	}
	if kafkaConfig.GetLogRetentionMs() != nil {
		res["log_retention_ms"] = strconv.FormatInt(kafkaConfig.GetLogRetentionMs().GetValue(), 10)
	}
	if kafkaConfig.GetLogSegmentBytes() != nil {
		res["log_segment_bytes"] = strconv.FormatInt(kafkaConfig.GetLogSegmentBytes().GetValue(), 10)
	}
	// if kafkaConfig.GetLogPreallocate() != nil {
	// 	res["log_preallocate"] = kafkaConfig.GetLogPreallocate().GetValue()
	// }
	if kafkaConfig.GetSocketSendBufferBytes() != nil {
		res["socket_send_buffer_bytes"] = strconv.FormatInt(kafkaConfig.GetSocketSendBufferBytes().GetValue(), 10)
	}
	if kafkaConfig.GetSocketReceiveBufferBytes() != nil {
		res["socket_receive_buffer_bytes"] = strconv.FormatInt(kafkaConfig.GetSocketReceiveBufferBytes().GetValue(), 10)
	}
	if kafkaConfig.GetAutoCreateTopicsEnable() != nil {
		res["auto_create_topics_enable"] = kafkaConfig.GetAutoCreateTopicsEnable().GetValue()
	}
	if kafkaConfig.GetNumPartitions() != nil {
		res["num_partitions"] = strconv.FormatInt(kafkaConfig.GetNumPartitions().GetValue(), 10)
	}
	if kafkaConfig.GetDefaultReplicationFactor() != nil {
		res["default_replication_factor"] = strconv.FormatInt(kafkaConfig.GetDefaultReplicationFactor().GetValue(), 10)
	}
	if kafkaConfig.GetMessageMaxBytes() != nil {
		res["message_max_bytes"] = strconv.FormatInt(kafkaConfig.GetMessageMaxBytes().GetValue(), 10)
	}
	if kafkaConfig.GetReplicaFetchMaxBytes() != nil {
		res["replica_fetch_max_bytes"] = strconv.FormatInt(kafkaConfig.GetReplicaFetchMaxBytes().GetValue(), 10)
	}
	if kafkaConfig.GetSslCipherSuites() != nil {
		res["ssl_cipher_suites"] = convertStringArrToInterface(kafkaConfig.GetSslCipherSuites())
	}
	if kafkaConfig.GetOffsetsRetentionMinutes() != nil {
		res["offsets_retention_minutes"] = strconv.FormatInt(kafkaConfig.GetOffsetsRetentionMinutes().GetValue(), 10)
	}
	if kafkaConfig.GetSaslEnabledMechanisms() != nil {
		res["sasl_enabled_mechanisms"] = convertStringArrToInterface(convertSaslEnabledMechanismsToStrings(kafkaConfig.GetSaslEnabledMechanisms()))
	}
	return res, nil
}

func convertSaslEnabledMechanismsToStrings(mechanisms []kafka.SaslMechanism) []string {
	var result []string
	for _, mechanism := range mechanisms {
		result = append(result, mechanism.String())
	}
	return result
}

func flattenKafkaConfig2_8Settings(r *kafka.KafkaConfig2_8) (map[string]interface{}, error) {
	return flattenKafkaConfigSettings(r)
}

func flattenKafkaConfig3Settings(r *kafka.KafkaConfig3) (map[string]interface{}, error) {
	return flattenKafkaConfigSettings(r)
}

func flattenKafkaResources(r *kafka.Resources) (map[string]interface{}, error) {
	res := map[string]interface{}{}

	res["resource_preset_id"] = r.ResourcePresetId
	res["disk_type_id"] = r.DiskTypeId
	res["disk_size"] = toGigabytes(r.DiskSize)

	return res, nil
}

func expandKafkaResources(d *schema.ResourceData, rootKey string) *kafka.Resources {
	resources := &kafka.Resources{}

	if v, ok := d.GetOk(rootKey + ".resource_preset_id"); ok {
		resources.ResourcePresetId = v.(string)
	}
	if v, ok := d.GetOk(rootKey + ".disk_size"); ok {
		resources.DiskSize = toBytes(v.(int))
	}
	if v, ok := d.GetOk(rootKey + ".disk_type_id"); ok {
		resources.DiskTypeId = v.(string)
	}
	return resources
}

func kafkaUserHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if n, ok := m["name"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", n.(string)))
	}
	if p, ok := m["password"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", p.(string)))
	}
	if ps, ok := m["permission"]; ok {
		permissions, _ := expandKafkaPermissions(ps.(*schema.Set))
		buf.WriteString(fmt.Sprintf("%s-", UserPermissionsToStr(permissions)))
	}
	return hashcode.String(buf.String())
}

func kafkaUserPermissionHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if n, ok := m["topic_name"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", n.(string)))
	}
	if r, ok := m["role"]; ok {
		buf.WriteString(fmt.Sprintf("%v-", r))
	}
	if allowHosts, ok := m["allow_hosts"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", AllowHostsSortedToStr(parseKafkaPermissionAllowHosts(allowHosts))))
	}
	return hashcode.String(buf.String())
}

func AllowHostsSortedToStr(allowHosts []string) string {
	if len(allowHosts) == 0 {
		return ""
	}
	sort.Strings(allowHosts)
	return strings.Join(allowHosts, ",")
}

func sortPermissions(permissions []*kafka.Permission) {
	sort.Slice(permissions, func(i, j int) bool {
		permFirst := permissions[i]
		permSecond := permissions[j]
		return permFirst.TopicName < permSecond.TopicName ||
			(permFirst.TopicName == permSecond.TopicName && permFirst.Role.String() < permSecond.Role.String()) ||
			(permFirst.TopicName == permSecond.TopicName && permFirst.Role.String() == permSecond.Role.String() &&
				AllowHostsSortedToStr(permFirst.GetAllowHosts()) < AllowHostsSortedToStr(permSecond.GetAllowHosts()))
	})
}

func userPermissionToStr(permission *kafka.Permission) string {
	return fmt.Sprintf("%s:%s:[%s]", permission.TopicName, permission.Role.String(), AllowHostsSortedToStr(permission.GetAllowHosts()))
}

func UserPermissionsToStr(permissions []*kafka.Permission) string {
	sortPermissions(permissions)
	strPermissionsSlice := []string{}
	for _, permission := range permissions {
		strPermissionsSlice = append(strPermissionsSlice, userPermissionToStr(permission))
	}
	return strings.Join(strPermissionsSlice, ",")
}

func kafkaHostHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if n, ok := m["name"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", n.(string)))
	}
	return hashcode.String(buf.String())
}

func flattenKafkaTopics(topics []*kafka.Topic) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)

	for _, d := range topics {
		m := make(map[string]interface{})
		m["name"] = d.GetName()
		m["partitions"] = d.GetPartitions().GetValue()
		m["replication_factor"] = d.GetReplicationFactor().GetValue()
		var cfg map[string]interface{}
		if d.GetTopicConfig_2_8() != nil {
			cfg = flattenKafkaTopicConfig2_8(d.GetTopicConfig_2_8())
		}
		if d.GetTopicConfig_3() != nil {
			cfg = flattenKafkaTopicConfig3(d.GetTopicConfig_3())
		}
		if len(cfg) != 0 {
			m["topic_config"] = []map[string]interface{}{cfg}
		}
		result = append(result, m)
	}

	return result
}

type TopicConfigSpec interface {
	GetCompressionType() kafka.CompressionType
	GetDeleteRetentionMs() *wrappers.Int64Value
	GetFileDeleteDelayMs() *wrappers.Int64Value
	GetFlushMessages() *wrappers.Int64Value
	GetFlushMs() *wrappers.Int64Value
	GetMinCompactionLagMs() *wrappers.Int64Value
	GetRetentionBytes() *wrappers.Int64Value
	GetRetentionMs() *wrappers.Int64Value
	GetMaxMessageBytes() *wrappers.Int64Value
	GetMinInsyncReplicas() *wrappers.Int64Value
	GetSegmentBytes() *wrappers.Int64Value
	GetPreallocate() *wrappers.BoolValue
}

func flattenKafkaTopicConfig(topicConfig TopicConfigSpec) map[string]interface{} {
	result := make(map[string]interface{})

	if topicConfig.GetCompressionType() != kafka.CompressionType_COMPRESSION_TYPE_UNSPECIFIED {
		result["compression_type"] = topicConfig.GetCompressionType().String()
	}
	if topicConfig.GetDeleteRetentionMs() != nil {
		result["delete_retention_ms"] = strconv.FormatInt(topicConfig.GetDeleteRetentionMs().GetValue(), 10)
	}
	if topicConfig.GetFileDeleteDelayMs() != nil {
		result["file_delete_delay_ms"] = strconv.FormatInt(topicConfig.GetFileDeleteDelayMs().GetValue(), 10)
	}
	if topicConfig.GetFlushMessages() != nil {
		result["flush_messages"] = strconv.FormatInt(topicConfig.GetFlushMessages().GetValue(), 10)
	}
	if topicConfig.GetFlushMs() != nil {
		result["flush_ms"] = strconv.FormatInt(topicConfig.GetFlushMs().GetValue(), 10)
	}
	if topicConfig.GetMinCompactionLagMs() != nil {
		result["min_compaction_lag_ms"] = strconv.FormatInt(topicConfig.GetMinCompactionLagMs().GetValue(), 10)
	}
	if topicConfig.GetRetentionBytes() != nil {
		result["retention_bytes"] = strconv.FormatInt(topicConfig.GetRetentionBytes().GetValue(), 10)
	}
	if topicConfig.GetRetentionMs() != nil {
		result["retention_ms"] = strconv.FormatInt(topicConfig.GetRetentionMs().GetValue(), 10)
	}
	if topicConfig.GetMaxMessageBytes() != nil {
		result["max_message_bytes"] = strconv.FormatInt(topicConfig.GetMaxMessageBytes().GetValue(), 10)
	}
	if topicConfig.GetMinInsyncReplicas() != nil {
		result["min_insync_replicas"] = strconv.FormatInt(topicConfig.GetMinInsyncReplicas().GetValue(), 10)
	}
	if topicConfig.GetSegmentBytes() != nil {
		result["segment_bytes"] = strconv.FormatInt(topicConfig.GetSegmentBytes().GetValue(), 10)
	}
	if topicConfig.GetPreallocate() != nil {
		result["preallocate"] = topicConfig.GetPreallocate().GetValue()
	}
	return result
}

func flattenKafkaTopicConfig2_8(topicConfig *kafka.TopicConfig2_8) map[string]interface{} {
	result := flattenKafkaTopicConfig(topicConfig)

	if topicConfig.GetCleanupPolicy() != kafka.TopicConfig2_8_CLEANUP_POLICY_UNSPECIFIED {
		result["cleanup_policy"] = topicConfig.GetCleanupPolicy().String()
	}

	return result
}

func flattenKafkaTopicConfig3(topicConfig *kafka.TopicConfig3) map[string]interface{} {
	result := flattenKafkaTopicConfig(topicConfig)

	if topicConfig.GetCleanupPolicy() != kafka.TopicConfig3_CLEANUP_POLICY_UNSPECIFIED {
		result["cleanup_policy"] = topicConfig.GetCleanupPolicy().String()
	}

	return result
}

func flattenKafkaUserPermissions(user *kafka.User) *schema.Set {
	result := schema.NewSet(kafkaUserPermissionHash, nil)
	for _, perm := range user.Permissions {
		p := map[string]interface{}{}
		p["topic_name"] = perm.TopicName
		p["role"] = perm.Role.String()
		if len(perm.GetAllowHosts()) > 0 {
			p["allow_hosts"] = convertStringArrToSchemaSet(perm.GetAllowHosts())
		}
		result.Add(p)
	}
	return result
}

func flattenKafkaUsers(users []*kafka.User, passwords map[string]string) *schema.Set {
	result := schema.NewSet(kafkaUserHash, nil)
	for _, user := range users {
		u := map[string]interface{}{}
		u["name"] = user.Name
		u["permission"] = flattenKafkaUserPermissions(user)
		if p, ok := passwords[user.Name]; ok {
			u["password"] = p
		}
		result.Add(u)
	}
	return result
}

func flattenKafkaHosts(hosts []*kafka.Host) *schema.Set {
	result := schema.NewSet(kafkaHostHash, nil)

	for _, host := range hosts {
		h := map[string]interface{}{}
		h["name"] = host.Name
		h["zone_id"] = host.ZoneId
		h["role"] = host.Role.String()
		h["health"] = host.Health.String()
		h["subnet_id"] = host.SubnetId
		h["assign_public_ip"] = host.AssignPublicIp

		result.Add(h)
	}
	return result
}

func kafkaUsersPasswords(users []*kafka.UserSpec) map[string]string {
	result := map[string]string{}
	for _, u := range users {
		result[u.Name] = u.Password
	}
	return result
}

func kafkaUsersDiff(currUsers []*kafka.User, targetUsers []*kafka.UserSpec) ([]string, []*kafka.UserSpec) {
	m := map[string]bool{}
	toAdd := []*kafka.UserSpec{}
	toDelete := map[string]bool{}
	for _, user := range currUsers {
		toDelete[user.Name] = true
		m[user.Name] = true
	}

	for _, user := range targetUsers {
		delete(toDelete, user.Name)
		_, ok := m[user.Name]
		if !ok {
			toAdd = append(toAdd, user)
			continue
		}
	}

	toDel := []string{}
	for u := range toDelete {
		toDel = append(toDel, u)
	}

	return toDel, toAdd
}

type entityDiff struct {
	OldEntityKey string
	OldEntity    map[string]interface{}
	NewEntityKey string
	NewEntity    map[string]interface{}
}

func diffByEntityKey(d *schema.ResourceData, path, indexKey string) map[string]entityDiff {
	result := map[string]entityDiff{}
	for i := 0; ; i++ {
		key := fmt.Sprintf("%s.%d", path, i)
		oldEntityI, newEntityI := d.GetChange(key)
		empty := true

		oldEntity := oldEntityI.(map[string]interface{})
		oldEntityKey, ok := oldEntity[indexKey].(string)
		if ok {
			empty = false
			diffEntity := result[oldEntityKey]
			diffEntity.OldEntity = oldEntity
			diffEntity.OldEntityKey = key
			result[oldEntityKey] = diffEntity
		}

		newEntity := newEntityI.(map[string]interface{})
		newEntityKey, ok := newEntity[indexKey].(string)
		if ok {
			empty = false
			if newEntityKey != "" {
				diffEntity := result[newEntityKey]
				diffEntity.NewEntity = newEntity
				diffEntity.NewEntityKey = key
				result[newEntityKey] = diffEntity
			}
		}

		if empty {
			break
		}
	}
	return result
}

func kafkaMaintenanceWindowSchemaValidateFunc(v interface{}, k string) (s []string, es []error) {
	dayString := v.(string)
	day, ok := kafka.WeeklyMaintenanceWindow_WeekDay_value[dayString]
	if !ok || day == 0 {
		es = append(es, fmt.Errorf(`expected %s value should be one of ("MON", "TUE", "WED", "THU", "FRI", "SAT", "SUN"). Current value is %v`, k, v))
		return
	}

	return
}

func flattenKafkaMaintenanceWindow(mw *kafka.MaintenanceWindow) ([]interface{}, error) {
	maintenanceWindow := map[string]interface{}{}
	if mw != nil {
		switch p := mw.GetPolicy().(type) {
		case *kafka.MaintenanceWindow_Anytime:
			maintenanceWindow["type"] = "ANYTIME"
		case *kafka.MaintenanceWindow_WeeklyMaintenanceWindow:
			maintenanceWindow["type"] = "WEEKLY"
			maintenanceWindow["hour"] = p.WeeklyMaintenanceWindow.Hour
			maintenanceWindow["day"] = kafka.WeeklyMaintenanceWindow_WeekDay_name[int32(p.WeeklyMaintenanceWindow.GetDay())]
		default:
			return nil, fmt.Errorf("unsupported Kafka maintenance policy type")
		}
	}

	return []interface{}{maintenanceWindow}, nil
}

func expandKafkaMaintenanceWindow(d *schema.ResourceData) (*kafka.MaintenanceWindow, error) {
	if _, ok := d.GetOk("maintenance_window"); !ok {
		return nil, nil
	}

	out := &kafka.MaintenanceWindow{}
	typeMW, _ := d.GetOk("maintenance_window.0.type")
	if typeMW == "ANYTIME" {
		if hour, ok := d.GetOk("maintenance_window.0.hour"); ok && hour != "" {
			return nil, fmt.Errorf("hour should not be set, when using ANYTIME")
		}
		if day, ok := d.GetOk("maintenance_window.0.day"); ok && day != "" {
			return nil, fmt.Errorf("day should not be set, when using ANYTIME")
		}
		out.Policy = &kafka.MaintenanceWindow_Anytime{
			Anytime: &kafka.AnytimeMaintenanceWindow{},
		}
	} else if typeMW == "WEEKLY" {
		hourInterface, ok := d.GetOk("maintenance_window.0.hour")
		if !ok {
			return nil, fmt.Errorf("hour should be set when using WEEKLY maintenance")
		}
		hour := hourInterface.(int)

		dayString := d.Get("maintenance_window.0.day").(string)

		day, ok := kafka.WeeklyMaintenanceWindow_WeekDay_value[dayString]
		if !ok || day == 0 {
			return nil, fmt.Errorf(`day value should be one of ("MON", "TUE", "WED", "THU", "FRI", "SAT", "SUN")`)
		}

		out.Policy = &kafka.MaintenanceWindow_WeeklyMaintenanceWindow{
			WeeklyMaintenanceWindow: &kafka.WeeklyMaintenanceWindow{
				Hour: int64(hour),
				Day:  kafka.WeeklyMaintenanceWindow_WeekDay(day),
			},
		}
	} else {
		return nil, fmt.Errorf("maintenance_window.0.type should be ANYTIME or WEEKLY")
	}

	return out, nil
}

func expandKafkaAccess(d *schema.ResourceData) *kafka.Access {
	if _, ok := d.GetOkExists("config.0.access"); !ok {
		return nil
	}
	out := &kafka.Access{}

	if v, ok := d.GetOk("config.0.access.0.data_transfer"); ok {
		out.DataTransfer = v.(bool)
	}
	return out
}

func flattenKafkaAccess(c *kafka.ConfigSpec) []map[string]interface{} {
	out := map[string]interface{}{}
	if c != nil && c.GetAccess() != nil {
		out["data_transfer"] = c.GetAccess().GetDataTransfer()
	}
	return []map[string]interface{}{out}
}

func expandKafkaRestAPI(d *schema.ResourceData) *kafka.ConfigSpec_RestAPIConfig {
	if _, ok := d.GetOkExists("config.0.rest_api"); !ok {
		return nil
	}
	out := &kafka.ConfigSpec_RestAPIConfig{}

	if v, ok := d.GetOk("config.0.rest_api.0.enabled"); ok {
		out.Enabled = v.(bool)
	}
	return out
}

func flattenKafkaRestAPI(c *kafka.ConfigSpec) []map[string]interface{} {
	out := map[string]interface{}{}
	if c != nil && c.GetRestApiConfig() != nil {
		out["enabled"] = c.GetRestApiConfig().GetEnabled()
	}
	return []map[string]interface{}{out}
}

func flattenKafkaConnectorMirrormaker(mm *kafka.ConnectorConfigMirrorMaker) ([]map[string]interface{}, error) {
	config := map[string]interface{}{
		"topics":             mm.Topics,
		"replication_factor": mm.ReplicationFactor.GetValue(),
	}
	sourceCluster, err := flattenKafkaClusterConnection(mm.SourceCluster)
	if err != nil {
		return nil, err
	}
	targetCluster, err := flattenKafkaClusterConnection(mm.TargetCluster)
	if err != nil {
		return nil, err
	}
	config["source_cluster"] = []map[string]interface{}{sourceCluster}
	config["target_cluster"] = []map[string]interface{}{targetCluster}
	return []map[string]interface{}{config}, nil
}

func flattenKafkaClusterConnection(cc *kafka.ClusterConnection) (map[string]interface{}, error) {
	config := map[string]interface{}{
		"alias": cc.Alias,
	}
	switch cc.GetClusterConnection().(type) {
	case *kafka.ClusterConnection_ThisCluster:
		config["this_cluster"] = []interface{}{map[string]interface{}{}}
	case *kafka.ClusterConnection_ExternalCluster:
		config["external_cluster"] = []map[string]interface{}{flattenKafkaExternalClusterConnection(cc.GetExternalCluster())}
	default:
		return nil, fmt.Errorf("cluster connection type of mirrormaker's cluster with alias %q not specified", cc.Alias)
	}
	return config, nil
}

func flattenKafkaExternalClusterConnection(ecc *kafka.ExternalClusterConnection) map[string]interface{} {
	return map[string]interface{}{
		"bootstrap_servers": ecc.BootstrapServers,
		"sasl_username":     ecc.SaslUsername,
		"sasl_mechanism":    ecc.SaslMechanism,
		"security_protocol": ecc.SecurityProtocol,
	}
}

func flattenKafkaConnectorS3Sink(s3Sink *kafka.ConnectorConfigS3Sink) ([]map[string]interface{}, error) {
	config := map[string]interface{}{
		"topics":                s3Sink.Topics,
		"file_compression_type": s3Sink.FileCompressionType,
		"file_max_records":      s3Sink.FileMaxRecords.GetValue(),
	}
	s3Connection, err := flattenS3Connection(s3Sink.GetS3Connection())
	if err != nil {
		return nil, err
	}
	config["s3_connection"] = []map[string]interface{}{s3Connection}
	return []map[string]interface{}{config}, nil
}

func flattenS3Connection(s3Conn *kafka.S3Connection) (map[string]interface{}, error) {
	config := map[string]interface{}{
		"bucket_name": s3Conn.BucketName,
	}
	switch s3Conn.GetStorage().(type) {
	case *kafka.S3Connection_ExternalS3:
		config["external_s3"] = []map[string]interface{}{flattenExternalS3Storage(s3Conn.GetExternalS3())}
	default:
		return nil, fmt.Errorf("this s3 connection type of s3-sink connector is not supported by current version of terraform provider")
	}
	return config, nil
}

func flattenExternalS3Storage(externalS3 *kafka.ExternalS3Storage) map[string]interface{} {
	return map[string]interface{}{
		"access_key_id": externalS3.AccessKeyId,
		"endpoint":      externalS3.Endpoint,
		"region":        externalS3.Region,
	}
}
