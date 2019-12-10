import { Component, Input, ChangeDetectionStrategy } from '@angular/core';
import { IconDefinition } from '@fortawesome/pro-light-svg-icons';

@Component({
  selector: 'app-button',
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <button [attr.aria-label]="ariaLabel">
      <fa-icon class="icon" *ngIf="icon" [icon]="icon"></fa-icon>

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
        border: none;
        color: #808080;
        font: normal 400 1.3rem Lato, sans-serif;
        padding: 0;
      }
    `,
    `
      button > .icon {
        font-size: 2.2rem;
        padding: 0 3.2rem;
      }
    `,
  ],
})
export class ButtonComponent {

  @Input()
  icon: IconDefinition;

  @Input()
  ariaLabel: string;

}
