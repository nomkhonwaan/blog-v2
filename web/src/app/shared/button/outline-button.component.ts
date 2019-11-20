import { Component, ChangeDetectionStrategy } from '@angular/core';

import { ButtonComponent } from './button.component';

@Component({
  selector: 'app-outline-button',
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <button [attr.aria-label]="ariaLabel">
      <ng-content></ng-content>
    </button>
  `,
  styles: [
    `
      :host {
          align-items: center;
          cursor: pointer;
          display: inline-flex;
          height: 100%;
          justify-content: center;
      }
    `,
    `
      button {
          background: none;
          border: .1rem solid #b3b3b3;
          border-radius: 100vh;
          color: #4d4d4d;
          font: normal 300 1.6rem Lato, sans-serif;
          min-width: 12.8rem;
          min-height: 4rem;
          padding: 0;
      }
    `,
  ],
})
export class OutlineButtonComponent extends ButtonComponent { }
