import { animate, state, style, transition, trigger } from '@angular/animations';
import { Component, Input } from '@angular/core';

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
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.scss'],
})
export class UserNavbarComponent {

  /**
   * An authenticated user info object
   */
  @Input()
  userInfo: UserInfo;

  /**
   * Use to toggle menu pane for showing or hiding
   */
  showMenu = false;

  toggleMenu(): void {
    this.showMenu = !this.showMenu;
  }

}
