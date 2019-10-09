/**
 * Type refers to the post size which is different in each size
 */
export enum Type {
  /**
   * Single is a type that use to display only single item on the page
   */
  Single = 1,

  /**
   * Medium is a type that use to list items on the same category / tag with abstraction paragraph and featured image
   */
  Medium,

  /**
   * Thumbnail is a type that use to list items on the related, sub-category, sub-tag with featured image only
   */
  Thumbnail,
}
