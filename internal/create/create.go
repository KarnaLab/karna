package create

func Run(name, layerName, runtime *string, withLayers *bool) {
	generateTemplate(name, layerName, runtime, withLayers)
}
