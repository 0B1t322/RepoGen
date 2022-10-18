package queryexpression

type QueryExpression[T any] struct {
	Expression T
	Or         []QueryExpression[T]
	And        []QueryExpression[T]
}
