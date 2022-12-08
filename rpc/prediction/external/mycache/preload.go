package mycache

type Preload struct {
	Struct            string
	PreloadableFields map[string]bool
}

func (c *cache) CreatePreloadCache() {
	cachedPreload := []Preload{
		{Struct: "User", PreloadableFields: map[string]bool{
			"BankInformationList": true,
		}},
	}
	for _, preload := range cachedPreload {
		c.cache.Set(preload.Struct, preload.PreloadableFields)
	}
}

func (c *cache) IsPreloadable(modelName string, key string) bool {
	v, ok := c.cache.Get(modelName)
	if !ok {
		return false
	}
	preloadableFields, ok := v.(map[string]bool)
	if !ok {
		return false
	}
	return preloadableFields[key]
}
