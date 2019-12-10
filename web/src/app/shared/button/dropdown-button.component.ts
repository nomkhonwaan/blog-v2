import { trigger, style, state, transition, animate } from '@angular/animations';
import { Component, Input, Output, EventEmitter } from '@angular/core';
import { faChevronDown, IconDefinition } from '@fortawesome/fontawesome-free-solid';

import { ButtonComponent } from './button.component';

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
  template: `
    <button [attr.aria-label]="ariaLabel" (click)="toggleDropdown()">
      <span [style.margin-left]="'auto'"></span>
      <span>{{selectedItem.label}}</span>
      <fa-icon class="icon" [icon]="faChevronDown"></fa-icon>
    </button>

    <ul [@slideUpDown]="showDropdown ? true : false" class="dropdown _list-unstyled">
      <li class="dropdown-item" *ngFor="let item of items" (click)="onClick(item);toggleDropdown()">
        {{item.label}}
      </li>
    </ul>
  `,
  styles: [
    `
      :host {
        display: block;
        position: relative;
      }
    `,
    `
      button {
        align-items: center;
        background: #0288d1;
        border: none;
        border-radius: 100vh;
        color: #fff;
        display: flex;
        font: normal 400 1.3rem Lato, sans-serif;
        height: 3.2rem;
        padding: 0 1.6rem;
        width: 100%;
      }
    `,
    `
      button > .icon {
        font-size: 0.8rem;
        margin-left: auto;
      }
    `,
    `
      button > .icon:after {
        clear: both;
      }
    `,
    `
      .dropdown {
        background: #fff;
        border: 0.1rem solid #ececec;
        border-radius: 0.4rem;
        padding: 1.6rem 0 !important;
        position: absolute;
        right: 0;
        top: 4.8rem;
        width: 100%;
        z-index: 98;
      }
    `,
    `
      .dropdown:before {
        background: #fff;
        border-top: 0.1rem solid #ececec;
        border-left: 0.1rem solid #ececec;
        content: "";
        display: block;
        height: 1.6rem;
        position: absolute;
        right: 0.8rem;
        width: 1.6rem;
        top: -0.8rem;
        transform: rotate(45deg);
      }
    `,
    `
      .dropdown-item {
        color: #333;
        cursor: pointer;
        font: normal 400 1.6rem Lato, sans-serif;
        padding: 0.8rem 1.6rem;
      }
    `,
  ],
})
export class DropdownButtonComponent extends ButtonComponent {

  @Input()
  items: Array<{ label: string, value?: any }>;

  @Input()
  selectedItem: { label: string, value?: any };

  @Output()
  change: EventEmitter<any> = new EventEmitter(null);

  /**
    * Use to toggle menu pane for showing or hiding
    */
  showDropdown = false;

  faChevronDown: IconDefinition = faChevronDown;

  onClick(selectedItem: { label: string, value?: any }): void {
    this.change.emit(selectedItem);
  }

  toggleDropdown(): void {
    this.showDropdown = !this.showDropdown;
  }

}
