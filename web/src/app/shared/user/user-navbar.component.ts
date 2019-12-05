import { Component } from '@angular/core';

import { UserComponent } from './user.component';

@Component({
  selector: 'app-user-navbar',
  template: `
    <img class="picture" [attr.src]="userInfo.picture" />

    <ul class="nav _list-unstyled">
      <li>
        <a [routerLink]="['/admin', 'new-post']">Draft a new post</a>
      </li>
      <li>
        <a (click)="toggleShowDraftPosts()">Display my draft posts</a>
      </li>
      <li class="horizontal-separator"></li>
      <li>
        <a [routerLink]="['/admin', 'stats']">Stats</a>
      </li>
      <li class="horizontal-separator"></li>
      <li>
        <a [routerLink]="['/user']">Profile</a>
      </li>
      <li>
        <a [routerLink]="['/user', 'settings']">Settings</a>
      </li>
      <li>
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
      .nav {
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
      .nav:before {
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
      .nav > li {
        padding: 1.6rem 3.2rem;
      }
    `,
    `
      .nav > li > a {
        color: #333;
        font: normal 400 1.6rem Lato, sans-serif;
      }
    `,
    `
      .nav > .horizontal-separator {
        border-bottom: .1rem solid #ececec;
        margin-bottom: 1.6rem;
        padding: 1.6rem 0 0 0;
      }
    `
  ],
})
export class UserNavbarComponent extends UserComponent {

  toggleShowDraftPosts(): void {

  }

}
