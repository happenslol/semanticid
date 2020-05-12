package semanticid

type SemanticIDBuilder struct {
	namespace  string
	collection string
	idProvider IDProvider
	from       string
	validate   bool
}

func Builder() *SemanticIDBuilder {
	return &SemanticIDBuilder{
		namespace:  DefaultNamespace,
		collection: DefaultCollection,
		idProvider: DefaultIDProvider,
		from:       "",
		validate:   true,
	}
}

func (b *SemanticIDBuilder) WithNamespace(namespace string) *SemanticIDBuilder {
	b.namespace = namespace
	return b
}

func (b *SemanticIDBuilder) WithCollection(collection string) *SemanticIDBuilder {
	b.collection = collection
	return b
}

func (b *SemanticIDBuilder) WithIDProvider(idp IDProvider) *SemanticIDBuilder {
	b.idProvider = idp
	return b
}

func (b *SemanticIDBuilder) FromString(s string) *SemanticIDBuilder {
	b.from = s
	return b
}

func (b *SemanticIDBuilder) NoValidate() *SemanticIDBuilder {
	b.validate = false
	return b
}

func (b *SemanticIDBuilder) Build() (SemanticID, error) {
	if b.from != "" {
		return fromStringWithParams(b.from, b.idProvider, b.validate)
	} else {
		return newWithParams(b.namespace, b.collection, b.idProvider)
	}
}
