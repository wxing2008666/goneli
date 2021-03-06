package goneli

import (
	"fmt"
	"testing"

	"github.com/obsidiandynamics/libstdgo/check"
	"github.com/obsidiandynamics/libstdgo/scribe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultKafkaConsumerProvider(t *testing.T) {
	c := Config{}
	c.SetDefaults()

	cons, err := c.KafkaConsumerProvider(&KafkaConfigMap{
		"session.timeout.ms": 1000,
	})
	assert.Nil(t, cons)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "Required property")
	}
}

func TestDefaultKafkaProducerProvider(t *testing.T) {
	c := Config{}
	c.SetDefaults()

	prod, err := c.KafkaProducerProvider(&KafkaConfigMap{
		"foo": "bar",
	})
	assert.Nil(t, prod)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "No such configuration property")
	}
}

func TestGetString(t *testing.T) {
	assert.Equal(t, "some-default", getString("some-default", func() (string, error) { return "", check.ErrSimulated }))
	assert.Equal(t, "some-string", getString("some-default", func() (string, error) { return "some-string", nil }))
}

func TestConfigString(t *testing.T) {
	cfg := Config{}
	cfg.SetDefaults()
	assert.Contains(t, cfg.String(), "Config[")
}

func TestValidateConfig_valid(t *testing.T) {
	cfg := Config{
		KafkaConfig:           KafkaConfigMap{},
		LeaderTopic:           "leader-topic",
		LeaderGroupID:         "leader-group-d",
		KafkaConsumerProvider: StandardKafkaConsumerProvider(),
		Scribe:                scribe.New(scribe.StandardBinding()),
		Name:                  "name",
	}
	cfg.SetDefaults()
	assert.Nil(t, cfg.Validate())
}

func TestValidateConfig_invalidLimits(t *testing.T) {
	cfg := Config{
		KafkaConfig:           KafkaConfigMap{},
		LeaderTopic:           "leader-topic",
		LeaderGroupID:         "leader-group-id",
		KafkaConsumerProvider: StandardKafkaConsumerProvider(),
		Scribe:                scribe.New(scribe.StandardBinding()),
		Name:                  "name",
		MinPollInterval:       Duration(0),
	}
	cfg.SetDefaults()
	assert.NotNil(t, cfg.Validate())
}

func TestValidateConfig_invalidName(t *testing.T) {
	cfg := Config{
		Name: "some%thing",
	}
	cfg.SetDefaults()
	err := cfg.Validate()
	require.NotNil(t, err)
	require.Equal(t, "Name: must be in a valid format.", err.Error())
}

func TestValidateConfig_invalidLeaderTopic(t *testing.T) {
	cfg := Config{
		LeaderTopic: "topic$#",
	}
	cfg.SetDefaults()
	fmt.Println(cfg)
	err := cfg.Validate()
	require.NotNil(t, err)
	require.Equal(t, "LeaderTopic: must be in a valid format.", err.Error())
}

func TestSanitiseName(t *testing.T) {
	cfg := Config{
		Name: "some%thing",
	}
	cfg.SetDefaults()
	err := cfg.Validate()
	require.NotNil(t, err)
	require.Equal(t, err.Error(), "Name: must be in a valid format.")
}

func TestValidateConfig_default(t *testing.T) {
	cfg := Config{}
	cfg.SetDefaults()

	assert.Nil(t, cfg.Validate())
}
