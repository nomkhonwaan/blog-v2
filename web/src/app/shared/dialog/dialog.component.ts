import { Component, EventEmitter, Output, Input, HostBinding } from '@angular/core';
import { trigger, style, state, transition, animate } from '@angular/animations';

@Component({
  animations: [
    trigger('fadeInOut', [
      state('true', style({
        display: 'initial',
        opacity: .16,
      })),
      state('false', style({
        display: 'none',
        opacity: 0,
      })),
      transition('* => true', [
        style({ display: 'initial' }),
        animate('.4s ease-in-out'),
      ]),
      transition('true => false', [
        animate('.4s ease-in-out', style({ opacity: 0 })),
        style({ display: 'none' }),
      ]),
    ]),
  ],
  selector: 'app-dialog',
  template: `
    <ng-content></ng-content>
  `,
  styles: [
    `
      :host {
        background: #333;
        cursor: pointer;
        display: none;
        height: 100%;
        min-height: 100vh;
        opacity: .16;
        position: absolute;
        left: 25.6rem;
        width: 100%;
        z-index: 99;
      }
    `,
  ],
})
export class DialogComponent {

  @Input()
  @HostBinding('@fadeInOut')
  show: boolean;

  @Output()
  whenClose: EventEmitter<null> = new EventEmitter();

  onClick(): void {
    this.whenClose.emit(null);
  }

}
