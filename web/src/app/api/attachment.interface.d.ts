/**
 * An uploaded file on the storage server
 */
interface Attachment {
  /**
   * An uploaded file name
   */
  fileName: string;

  /**
   * Valid URL string composes with file name and ID
   */
  slug: string;
}
