import { Component, EventEmitter, Output } from '@angular/core';

@Component({
  selector: 'app-dialog',
  template: `
    <div click="onClick()">
      <ng-content></ng-content>
    </div>
  `,
  styles: [
    `
      :host {
          height: 100%;
          min-height: 100vh;
          position: absolute;
          left: 25.6rem;
          width: 100%;
      }
    `,
    `
      :host > div {
          background: #333;
          cursor: pointer;
          height: 100%;
          opacity: .16;
      }
    `,
  ],
})
export class DialogComponent {

  @Output()
  hide: EventEmitter<null> = new EventEmitter();

  onClick(): void {
    this.hide.emit(null);
  }

}
