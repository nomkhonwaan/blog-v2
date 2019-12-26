import { Location } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { ActivatedRoute, Router } from '@angular/router';
import { faArrowLeft, faBars, faSpinnerThird, IconDefinition } from '@fortawesome/pro-light-svg-icons';
import { select, Store } from '@ngrx/store';
import { Apollo } from 'apollo-angular';
import { ApolloQueryResult } from 'apollo-client';
import { GraphQLError } from 'graphql';
import gql from 'graphql-tag';
import { map, tap } from 'rxjs/operators';
import { toggleEditorSidebar } from 'src/app/app.actions';
import { slideInOut } from 'src/app/shared';
import { environment } from '../../../environments/environment';

@Component({
  animations: [
    slideInOut,
  ],
  selector: 'app-editor',
  templateUrl: './editor.component.html',
  styleUrls: ['./editor.component.scss'],
})
export class EditorComponent implements OnInit {

  /**
   * A post object
   */
  post: Post;

  /**
   * Use to display when GraphQL returns errors
   */
  errors: ReadonlyArray<GraphQLError> = null;

  /**
 * Use to toggle application sidebar
 */
  hasSidebarExpanded = false;

  /**
   * Use to indicate loading status
   */
  isFetching = false;

  /**
   * An authenticated user info object
   */
  userInfo: UserInfo | null;

  /**
   * To-be updated attachment
   */
  selectedAttachment: Attachment;

  /**
   * List of FontAwesome icons
   */
  icons: { [name: string]: IconDefinition } = {
    faArrowLeft,
    faBars,
    faSpinnerThird,
  };

  constructor(
    private apollo: Apollo,
    private location: Location,
    private route: ActivatedRoute,
    private router: Router,
    private store: Store<{ app: AppState }>,
    private title: Title,
  ) { }

  ngOnInit(): void {
    this.store.pipe(select('app')).subscribe((app: AppState): void => {
      this.hasSidebarExpanded = !app.admin.editor.sidebar.collapsed;
      this.userInfo = app.auth.userInfo;
    });

    const slug: string | null = this.route.snapshot.paramMap.get('slug');

    if (slug) {
      this.findPostBySlug(slug);
    } else {
      this.createNewPost();
    }
  }

  goBack(): void {
    this.location.back();
  }

  toggleSidebar(): void {
    this.store.dispatch(toggleEditorSidebar());
  }

  onChaging(isFetching: boolean): void {
    this.isFetching = isFetching;
  }

  onChangingSuccess(post: Post): void {
    this.post.slug = post.slug;
    this.post.status = post.status;
    this.post.html = post.html;
    this.post.publishedAt = post.publishedAt;
    this.post.categories = post.categories;
    this.post.tags = post.tags;
    this.post.featuredImage = post.featuredImage;
    this.post.attachments = post.attachments;
    this.post.updatedAt = post.updatedAt;

    this.title.setTitle(`Edit · ${post.title} - ${environment.title}`);
  }

  onChagningErrors(errors: ReadonlyArray<GraphQLError>): void {
    this.errors = errors;
  }

  onSelectingAttachment(attachment: Attachment): void {
    this.selectedAttachment = attachment;
  }

  private createNewPost(): void {
    this.title.setTitle(`Draft a new post - ${environment.title}`);

    this.apollo.mutate({
      mutation: gql`
        mutation {
          createPost {
            slug createdAt
          }
        }
      `,
    }).pipe(
      map((result: ApolloQueryResult<{ createPost: Post }>): Post => result.data.createPost),
    ).subscribe((post: Post): void => {
      const createdAt: Date = new Date(
        new Date(post.createdAt)
          .toLocaleString('en-US', { timeZone: 'Asia/Bangkok' }),
      );

      this.router.navigate([
        createdAt.getFullYear().toString(),
        (createdAt.getMonth() + 1).toString(),
        createdAt.getDate().toString(),
        post.slug,
        'edit',
      ]);
    });
  }

  private findPostBySlug(slug: string): void {
    this.apollo.query({
      query: gql`
        {
          post(slug: $slug) {
            title slug
            status
            markdown html
            publishedAt
            authorId
            categories { name slug }
            tags { name slug }
            featuredImage { slug }
            attachments { fileName slug }
            createdAt updatedAt
          }
        }
      `,
      variables: {
        slug,
      },
    }).pipe(
      map((result: ApolloQueryResult<{ post: Post }>): Post => result.data.post),
      tap((post: Post): void => {
        this.title.setTitle(`Edit · ${post.title || 'Untitled'} - ${environment.title}`);
      }),
    ).subscribe((post: Post): void => {
      this.post = post;
    });
  }
}
