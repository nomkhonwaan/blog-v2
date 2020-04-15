/**
 * A label attached to the post for the purpose of identification
 */
interface Tag {
  /**
   * Name of the tag
   */
  name: string;

  /**
   * Valid URL string composes with name and ID
   */
  slug: string;

  /**
   * List of latest published posts are belongging to the tag
   */
  latestPublishedPosts: Post[];
}
