# SemanticID

[![Build Status](https://travis-ci.org/gin-gonic/gin.svg)](https://travis-ci.org/happenslol/semanticid)
[![codecov](https://codecov.io/gh/happenslol/semanticid/branch/master/graph/badge.svg)](https://codecov.io/gh/happenslol/semanticid)
[![GoDoc](https://godoc.org/github.com/happenslol/semanticid?status.svg)](https://godoc.org/github.com/happenslol/semanticid)
[![Go Report Card](https://goreportcard.com/badge/github.com/happenslol/semanticid)](https://goreportcard.com/report/github.com/happenslol/semanticid)

SemanticIDs are a special type of ID, providing extra utility especially in the context of microservice infrastructures.

SemanticIDs consist of 3 parts:

```
<namespace>.<collection>.<id>
```

By choosing a meaningful namespace and collection, you can glean a lot of information about an ID just by looking at it. The unique ID part uses [ULIDs](<https://github.com/ulid/spec>) by default, but can easily be changed to using [UUIDv4](<https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)>).

You can also optionally add in your own ID provider, which will allow you to customize ID generation as well as validation. This can be done globally or on a case-by-case basis by using the Builder.

## Usage

SemanticID uses go modules internally, so it will seamlessly integrate with other projects using modules. This also means that **go 1.11+ is required**.  
To use the library, simply do:

```bash
$ go get -u github.com/happenslol/semanticid
```

Then, import it into your code:

```go
import "github.com/happenslol/semanticid"
```

Here's a simple example that shows how to create and parse IDs:

```go
semanticid.DefaultNamespace = "myservice"
semanticid.DefaultCollection = "entities"

sid := semanticid.Must(semanticid.NewDefault())

parsed, err := semanticid.FromString(toParse)
if err != nil {
  fmt.Printf("something went wrong while parsing: %v", err)
}
```

Here are some more examples for common use cases:

```
// Switch to UUIDv4
semanticid.DefaultIDProvider = semanticid.NewUUIDProvider()

// Use a custom provider just for a single ID
type MyProvider struct {}

func (p *MyProvider) Generate() (string, error) {
  return "1234", nil
}

func (p *MyProvider) Validate(id string) {
  return nil
}

customID, err := semanticid.Builder().
  WithIDProvider(&MyProvider{}).
  WithCollection("custom-id-entity").
  Build()

// Parse a semantic id without ID validation
sid, _ := semanticid.Builder().
  FromString("a.b.c").
  NoValidate().
  Build()
```

## Choosing namespace and collection

While you can generally choose any namespace and collection you want, here are a few guidelines that should make SemanticIDs more useful and consistent throughout your infrastructure:

**Choose a namespace that suggests in which part of your infrastructure the ID was created.** This could be the name of the microservice, the name of the external service the ID refers to, or your company's or project's name.

**Use the plural of the entity name as the collection.** This is just a convention, but it generally leads to very readable IDs and connects well with other parts of your applications - in practice, you can probably reuse that database name or URL for a specific entity as the collection name.

**Only use lowercase letters and no special characters in the namespace and collection.** This reduces visual noise and makes sure your IDs always stay URL safe and unambiguous.

Examples for good SemanticIDs:

```
accountservice.users.7da57b46-f4f4-4824-a8e8-0c05ff88d9a5
```

```
github.repos.87961165-15f0-4fb8-8d8b-d9ce59034565
```

```
blog.posts.59731722-54ea-4447-8e99-f4689c0c060a
```
