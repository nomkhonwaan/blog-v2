import { Component, OnInit, Input, HostBinding, ChangeDetectionStrategy } from '@angular/core';

import { PostTitleComponent } from './post-title.component';

@Component({
  selector: 'app-post-featured-image',
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <a [routerLink]="isDisabledLink() ? null : href" [attr.aria-label]="post.title">
      <img *ngIf="src" [src]="src" [alt]="post.title" class="lazyload">
    </a>
  `,
  styles: [
    `
      :host {
        display: block;
      }
    `,
    `
      img {
        border-radius: .4rem;
        max-width: 100%;
      }
    `
  ],
})
export class PostFeaturedImageComponent extends PostTitleComponent implements OnInit {

  @HostBinding('class.-with-featured-image')
  withFeaturedImage = false;

  @HostBinding('style.background-image')
  src: string;

  ngOnInit(): void {
    super.ngOnInit();

    if (this.hasFeaturedImage()) {
      this.src = `/api/v2/storage/${this.post.featuredImage.slug}?width=${this.innerWidth}&height=${this.innerHeight}`;
      this.withFeaturedImage = true;
    }
  }

  hasFeaturedImage(): boolean {
    return !!this.post && this.post.featuredImage.slug.length > 0;
  }

}
