import { Component, Input, Output, EventEmitter } from '@angular/core';

@Component({
  selector: 'app-attachment-viewer',
  template: `
    <app-dialog *ngIf="attachment" state="show">
      <div class="_flex _flex-horizontal-align-right _flex-vertical-align-middle" [style.margin]="'3.2rem 3.2rem 0 0'">
        <app-button [style.margin-right]="'-1.6rem'" (click)="onClick()" [style.margin-top]="'.2rem'">
          {{'close' | uppercase}}
        </app-button>
        <app-button ariaLabel="Close Sidebar" [icon]="faTimes" (click)="onClick()"></app-button>
      </div>

      <div class="viewer">
        <img src="/api/v2.1/storage/{{selectedAttachment.slug}}?width=420" alt="{{selectedAttachment.fileName}}" />

        <div [style.margin-top]="'2.4rem'">
          <input class="attachment-url" value="/api/v2.1/storage/{{selectedAttachment.slug}}" />
        </div>
      </div>
    </app-dialog>
  `,
  styles: [
    `
      app-dialog {
        background: rgba(51, 51, 51, 0.8);
        bottom: 0;
        left: 0;
        opacity: 1 !important;
        right: 0;
        top: 0;
      }
    `,
    `
      :ng-deep app-button > button {
        color: #fff!important;
      }
    `,
    `
      .viewer {
        background: #fff;
        border-radius: 0.2rem;
        display: inline-block;
        left: 50%;
        max-width: 48rem;
        padding: 3.2rem;
        position: absolute;
        top: 50%;
        transform: translate(-50%, -50%);
      }
    `,
    `
      .viewer > img {
        max-height: 42rem;
        max-width: 100%;
      }
    `,
    `
      .attachment-uri {
        border: 0.1rem solid #ececec;
        border-radius: 0.2rem;
        font: normal 400 1.6rem Lato, sans-serif;
        height: 3.2rem;
        line-height: 3.2rem;
        width: 100%;
      }
    `,
  ],
})
export class AttachmentViewerComponent {

  @Input()
  attachment: Attachment;

  @Output()
  click: EventEmitter<null> = new EventEmitter(null);

  onClick(): void {
    this.click.emit(null);
  }

}
