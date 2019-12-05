import { Component } from '@angular/core';

import { UserComponent } from './user.component';
import { trigger, style, state, transition, animate } from '@angular/animations';

@Component({
  animations: [
    trigger('slideUpDown', [
      state('true', style({
        transform: 'translateY(0)',
        opacity: 1,
        display: 'initial',
      })),
      state('false', style({
        transform: 'translateY(8%)',
        opacity: 0,
        display: 'none',
      })),
      transition('* => true', [
        style({ display: 'initial' }),
        animate('.4s ease-in-out', style({
          transform: 'translateY(0)',
          opacity: 1,
        })),
      ]),
      transition('true => false', [
        animate('.4s ease-in-out', style({
          transform: 'translateY(8%)',
          opacity: 0,
        })),
        style({ display: 'none' }),
      ]),
    ]),
  ],
  selector: 'app-user-navbar',
  template: `
    <img
      class="picture"
      [attr.src]="userInfo.picture"
      (click)="toggleMenu()"
    />

    <ul [@slideUpDown]="showMenu ? true : false" class="menu _list-unstyled">
      <li class="menu-item">
        <a [routerLink]="['/admin', 'new-post']">Draft a new post</a>
      </li>
      <li class="menu-item">
        <a [routerLink]="['/admin', 'posts']">All my posts</a>
      </li>
      <li class="menu-item menu-item--horizontal-separator"></li>
      <li class="menu-item">
        <a [routerLink]="['/admin', 'stats']">Stats</a>
      </li>
      <li class="menu-item menu-item--horizontal-separator"></li>
      <li class="menu-item">
        <a [routerLink]="['/user']">Profile</a>
      </li>
      <li class="menu-item">
        <a [routerLink]="['/user', 'settings']">Settings</a>
      </li>
      <li class="menu-item">
        <a [routerLink]="['/logout']">Logout</a>
      </li>
    </ul>
  `,
  styles: [
    `
      :host {
        position: relative;
      }
    `,
    `
      .picture {
        border-radius: 50%;
        cursor: pointer;
        max-height: 3.2rem;
        max-width: 3.2rem;
      }
    `,
    `
      .menu {
        background: #fff;
        border: .1rem solid #ececec;
        border-radius: .4rem;
        min-width: 24rem;
        padding: 1.6rem 0!important;
        position: absolute;
        right: 0;
        top: 7.2rem;
      }
    `,
    `
      .menu:before {
        background: #fff;
        border-top: .1rem solid #ececec;
        border-left: .1rem solid #ececec;
        content: '';
        display: block;
        height: 1.6rem;
        position: absolute;
        right: .8rem;
        width: 1.6rem;
        top: -.8rem;
        transform: rotate(45deg);
      }
    `,
    `
      .menu-item {
        padding: 1.6rem 3.2rem;
      }
    `,
    `
      .menu-item > a {
        color: #333;
        font: normal 400 1.6rem Lato, sans-serif;
      }
    `,
    `
      .menu-item--horizontal-separator {
        border-bottom: .1rem solid #ececec;
        margin-bottom: 1.6rem;
        padding: 1.6rem 0 0 0;
      }
    `
  ],
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
