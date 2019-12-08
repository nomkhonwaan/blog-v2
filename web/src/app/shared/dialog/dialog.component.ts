import { trigger, state, style, transition, animate } from '@angular/animations';
import { Component, Input, HostBinding } from '@angular/core';

@Component({
  animations: [
    trigger('fadeInOut', [
      state('show', style({
        display: 'block',
        opacity: .16,
      })),
      state('hide', style({
        display: 'none',
        opacity: 0,
      })),
      // transition('* => show', [
      //   style({ display: 'block', opacity: 0 }),
      //   animate('.4s ease-in-out', style({ opacity: .16 })),
      // ]),
      // transition('show => hide', [
      //   animate('.4s ease-in-out', style({ opacity: 0 })),
      //   style({ display: 'none' }),
      // ]),
    ]),
  ],
  selector: 'app-dialog',
  template: `
    <ng-content></ng-content>
  `,
  styleUrls: ['./dialog.component.scss'],
})
export class DialogComponent {

  @Input()
  @HostBinding('@fadeInOut')
  state: string = 'hide';

}
