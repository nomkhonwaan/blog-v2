import { state, style, trigger } from '@angular/animations';
import { Component, HostBinding, Input } from '@angular/core';

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
  templateUrl: './dialog.component.html',
  styleUrls: ['./dialog.component.scss'],
})
export class DialogComponent {

  /**
   * Use to indiciate whether dialog should show or not
   */
  @Input()
  @HostBinding('@fadeInOut')
  state = 'hide';

}
