package graphql

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/podops/podops"
	"github.com/podops/podops/graphql/graph"
	"github.com/podops/podops/graphql/graph/generated"
)

func TestGraphqlSchema(t *testing.T) {

	resolver := graph.CreateResolver()
	assert.NotNil(t, resolver)

	schema := generated.NewExecutableSchema(generated.Config{Resolvers: resolver})
	assert.NotNil(t, schema)
}

func TestGraphqlShow(t *testing.T) {
	ctx := context.TODO()

	opts := podops.LoadConfiguration()
	assert.NotNil(t, opts)

	client, err := podops.NewClient(ctx, opts.Token)
	if assert.NoError(t, err) {
		assert.NotNil(t, client)

		prod, err := client.Productions()
		if assert.NoError(t, err) {
			assert.NotNil(t, prod)
			assert.GreaterOrEqual(t, len(prod.Productions), 1)

			name := prod.Productions[0].Name

			resolver := graph.CreateResolver()
			show, err := resolver.Query().Show(ctx, &name, 5)

			if assert.NoError(t, err) {
				assert.NotNil(t, show)
				assert.Equal(t, show.Name, name)
			}
		}
	}
}

func TestGraphqlEpisode(t *testing.T) {
	ctx := context.TODO()

	opts := podops.LoadConfiguration()
	assert.NotNil(t, opts)

	client, err := podops.NewClient(ctx, opts.Token)
	if assert.NoError(t, err) {
		assert.NotNil(t, client)

		prod, err := client.Productions()
		if assert.NoError(t, err) {

			episodes, err := client.Resources(prod.Productions[0].GUID, podops.ResourceEpisode)
			if assert.NoError(t, err) {
				assert.NotNil(t, episodes)
				assert.GreaterOrEqual(t, len(episodes.Resources), 1)

				guid := episodes.Resources[0].GUID

				resolver := graph.CreateResolver()
				episode, err := resolver.Query().Episode(ctx, &guid)

				if assert.NoError(t, err) {
					assert.NotNil(t, episode)
					assert.Equal(t, episodes.Resources[0].GUID, episode.GUID)
				}
			}
		}
	}
}
