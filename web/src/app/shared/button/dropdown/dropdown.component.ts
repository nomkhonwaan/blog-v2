import { animate, state, style, transition, trigger } from '@angular/animations';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { faChevronDown, IconDefinition } from '@fortawesome/fontawesome-free-solid';
import { ButtonComponent } from '../button';

@Component({
  animations: [
    trigger('slideUpDown', [
      state('true', style({
        display: 'block',
        opacity: 1,
        transform: 'translateY(0)',
      })),
      state('false', style({
        display: 'none',
        opacity: 0,
        transform: 'translateY(8%)',
      })),
      transition('* => true', [
        style({ display: 'block' }),
        animate('.4s ease-in-out', style({
          opacity: 1,
          transform: 'translateY(0)',
        })),
      ]),
      transition('true => false', [
        animate('.4s ease-in-out', style({
          opacity: 0,
          transform: 'translateY(8%)',
        })),
        style({ display: 'none' }),
      ]),
    ]),
  ],
  selector: 'app-dropdown-button',
  templateUrl: './dropdown.component.html',
  styleUrls: ['./dropdown.component.scss'],
})
export class DropdownButtonComponent extends ButtonComponent {

  /**
   * List of items to-be displayed
   */
  @Input()
  items: Array<DropdownItem>;

  /**
   * An item that has been selected from the list of items
   */
  @Input()
  selectedItem: DropdownItem;

  /**
   * For emitting a selected item
   */
  @Output()
  selectItem: EventEmitter<DropdownItem> = new EventEmitter(null);

  /**
   * Use to indicate whether dropdown menu should show or not
   */
  showMenu = false;

  /**
   * A FontAwesome icon to-be displayed with the button
   */
  faChevronDown: IconDefinition = faChevronDown;

  onClick(selectedItem: { label: string, value?: any }): void {
    this.selectItem.emit(selectedItem);
  }

  toggleMenu(): void {
    this.showMenu = !this.showMenu;
  }

}
