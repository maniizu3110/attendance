package mycache

type Joins struct {
	Struct         string
	JoinableFields map[string]bool
}

// example
func (c *cache) CreateJoinsCache() {
	cachedJoins := []Joins{
		{Struct: "Music", JoinableFields: map[string]bool{
			"Split": true,
		}},
	}
	for _, joins := range cachedJoins {
		c.cache.Set(joins.Struct, joins.JoinableFields)
	}
}

func (c *cache) IsJoinable(modelName string, key string) bool {
	v, ok := c.cache.Get(modelName)
	if !ok {
		return false
	}
	JoinableFields, ok := v.(map[string]bool)
	if !ok {
		return false
	}
	return JoinableFields[key]
}
