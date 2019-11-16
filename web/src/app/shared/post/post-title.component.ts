import { Component, OnInit, Input, ChangeDetectionStrategy } from '@angular/core';

import { PostComponent } from './post.component';

@Component({
  selector: 'app-post-title',
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <a [routerLink]="isDisabledLink() ? null : href">
      <ng-content></ng-content>
    </a>
  `,
  styles: [
    `
      :host {
        color: #333;
        display: block;
        font-family: Athiti, sans-serif;
        font-weight: 500;
      }
    `,
    `
      ::ng-deep h1 {
        font-size: 4.8rem;
        font-weight: inherit;
      }
    `,
    `
      ::ng-deep h2 {
        font-size: 3.6rem;
        font-weight: inherit;
      }
    `,
  ],
})
export class PostTitleComponent extends PostComponent implements OnInit {

  /**
   * Used to indiciate whether link or no-link
   */
  @Input()
  link = true;

  href: string[];

  ngOnInit(): void {
    const publishedAt: Date = new Date(
      new Date('2016-04-01T17:00:00Z')
        .toLocaleString('en-US', { timeZone: 'Asia/Bangkok' }),
    );

    if (this.link) {
      this.href = [
        publishedAt.getFullYear().toString(),
        (publishedAt.getMonth() + 1).toString(), // a month number start at 0 not 1
        publishedAt.getDate().toString(),
        this.post.slug,
      ];
    }
  }

  isDisabledLink(): boolean {
    return !!!this.link;
  }

}
