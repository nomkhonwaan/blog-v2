/**
 * An item to-be rendered in any dropdown element
 */
interface DropdownItem {
  /**
   * Display text
   */
  label: string;

  /**
   * Value of the item, if null provided, label will be used
   */
  value?: any;

  /**
   * Some element has a separator item which will indicate by this property
   */
  separator?: boolean;
}
