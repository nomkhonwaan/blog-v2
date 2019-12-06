import { trigger, style, state, transition, animate } from '@angular/animations';
import { Component } from '@angular/core';

import { UserComponent } from './user.component';

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
  selector: 'app-user-navbar',
  templateUrl: './user-navbar.component.html',
  styleUrls: ['./user-navbar.component.scss'],
})
export class UserNavbarComponent extends UserComponent {

  /**
   * Used to toggle menu pane for showing or hiding
   */
  showMenu = false;

  toggleMenu(): void {
    this.showMenu = !this.showMenu;
  }

}
