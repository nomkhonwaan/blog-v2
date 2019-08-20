// A piece of content in the blog platform
interface Post {
  // Identifier of the post
  id: string

  // Title of the post
  title: string

  // Valid URL string composes with title and ID
  slug: string

  // Status of the post which could be...
  // - PUBLISHED
  // - DRAFT
  status: string

  // Original content of the post in markdown syntax
  markdown: string

  // Content of the post in HTML format which will be translated from markdown
  html: string

  // Date-time that the post was published
  publishedAt: Date

  // Identifier of the author
  authorId: string

  // Date-time that the post was created
  createdAt: Date

  // Date-time that the post was updated
  updatedAt: Date
}
