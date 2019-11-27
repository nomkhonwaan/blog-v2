/**
 * A piece of content in the blog platform
 */
interface Post {
  /**
   * Title of the post
   */
  title: string;

  /**
   * Valid URL string composes with title and ID
   */
  slug: string;

  /**
   * Status of the post which could be...
   * - PUBLISHED
   * - DRAFT
   */
  status: string;

  /**
   * Original content of the post in markdown syntax
   */
  markdown: string;

  /**
   * Content of the post in HTML format which will be translated from markdown
   */
  html: string;

  /**
   * Date-time that the post was published
   */
  publishedAt: Date;

  /**
   * Identifier of the author
   */
  authorId: string;

  /**
   * List of categories that the post belongging to
   */
  categories: Category[];

  /**
   * List of tags that the post belongging to
   */
  tags: Tag[];

  /**
   * A featured image to be shown in the social network as a cover image
   */
  featuredImage: File;

  /**
   * List of attachments are belonging to the post
   */
  Attachments: File[];

  /**
   * A social network engagement of the post
   */
  engagement: Engagement;

  /**
   * Date-time that the post was created
   */
  createdAt: Date;

  /**
   * Date-time that the post was updated
   */
  updatedAt: Date;
}
